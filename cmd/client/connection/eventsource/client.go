package eventsource

import (
	"bufio"
	"bytes"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ecromaneli-golang/http/webserver"
)

type EventSource struct {
	url          string
	header       http.Header
	lastEventId  string
	retryTimeout int
	request      *http.Request
	response     *http.Response
	stream       chan *webserver.Event
	err          chan error
	closed       chan any
}

var _DEFAULT_REQUEST_HEADER_ = map[string]string{
	"Accept": "text/event-stream",
	// "Cache-Control": "no-cache",
	// "Connection":    "keep-alive",
}

func (this *EventSource) Listen() <-chan *webserver.Event {
	return this.stream
}

func (this *EventSource) Error() <-chan error {
	return this.err
}

func (this *EventSource) sendErrorIfNotNil(err error) {
	if err != nil {
		select {
		case this.err <- err:
		default:
		}
	}
}

func (this *EventSource) Closed() bool {
	select {
	case <-this.closed:
		return true
	default:
		return false
	}
}

func (this *EventSource) Close() {
	if !this.Closed() {
		close(this.closed)
	}
}

func New(url string) (*EventSource, error) {
	return NewWithHeader(url, nil)
}

func NewWithHeader(url string, header http.Header) (*EventSource, error) {
	client := &EventSource{
		url:    url,
		header: header,
		stream: make(chan *webserver.Event, 32),
		closed: make(chan any),
	}

	go func() {
		for !client.Closed() {
			statusCode, err := client.connect()

			client.sendErrorIfNotNil(err)

			// 204 represents that no more content will be sent
			if statusCode == http.StatusNoContent {
				client.Close()
				continue
			}

			client.startListening()
			<-time.After(time.Duration(client.retryTimeout) * time.Millisecond)
		}
	}()

	return client, nil
}

func (this *EventSource) connect() (statusCode int, err error) {
	req, err := http.NewRequest("GET", this.url, nil)

	if err != nil {
		return http.StatusBadRequest, err
	}

	this.request = req

	for key, value := range _DEFAULT_REQUEST_HEADER_ {
		req.Header.Set(key, value)
	}

	if this.header != nil {
		for key, values := range this.header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}

	if this.lastEventId != "" {
		req.Header.Set("Last-Event-ID", this.lastEventId)
	}

	// Avoid TCP reuse and do request
	req.Close = true
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return http.StatusBadRequest, err
	}

	this.response = resp

	if strings.Index(resp.Header.Get("Content-Type"), _DEFAULT_REQUEST_HEADER_["Accept"]) == -1 {
		return resp.StatusCode, errors.New("Content-Type not accepted")
	}

	return resp.StatusCode, nil
}

func (this *EventSource) startListening() {
	defer this.response.Body.Close()
	this.sendErrorIfNotNil(this.scanEvents())
}

func (this *EventSource) scanEvents() error {
	scanner := bufio.NewScanner(this.response.Body)
	scanner.Split(scanLines)

	event := &webserver.Event{}

	for scanner.Scan() {
		line := scanner.Bytes()

		if len(line) == 0 {
			this.stream <- event
			event = &webserver.Event{}
			continue
		}

		key, value := tokenize(line)

		switch string(key) {
		case "event":
			event.Name = string(value)
		case "data":
			event.Data = string(value)
		case "id":
			event.ID = string(value)
			this.lastEventId = event.ID
		case "retry":
			retry, err := strconv.Atoi(string(value))

			if err != nil {
				return err
			}

			if retry < 1 {
				return errors.New("Invalid retry timeout '" + string(value) + "'")
			}

			this.retryTimeout = retry

		default:
			return errors.New("Invalid token '" + key + "'")
		}
	}

	return scanner.Err()
}

func tokenize(line []byte) (string, []byte) {
	var key []byte
	var value []byte

	colon := bytes.IndexByte(line, ':')

	if colon == -1 {
		return string(line), []byte{}
	}

	if colon == 0 {
		return "", nil
	}

	key = line[:colon]
	value = line[colon+1:]

	if value[0] == ' ' {
		value = value[1:]
	}

	return string(key), value
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

// ScanLines is a split function for a Scanner that returns each line of
// text, stripped of any trailing end-of-line marker. The returned line may
// be empty. The end-of-line marker is one optional carriage return followed
// by one mandatory newline. In regular expression notation, it is `\r?\n`.
// The last non-empty line of input will be returned even if it has no
// newline.
func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

func panicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}
