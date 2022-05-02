package eventsource

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"woole/app"
)

type Client struct {
	Stream <-chan Event
	Err    error
}

var _DEFAULT_REQUEST_HEADER_ = map[string]string{
	"Accept": "text/event-stream",
	// "Cache-Control": "no-cache",
	// "Connection":    "keep-alive",
}

func NewRequest(eventsourceUrl string) (*Client, error) {
	req, err := http.NewRequest("GET", eventsourceUrl, nil)

	if err != nil {
		return nil, err
	}

	for key, value := range _DEFAULT_REQUEST_HEADER_ {
		req.Header.Set(key, value)
	}

	app.SetAuthorization(req.Header)

	// Avoid TCP reuse and do request
	req.Close = true
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	// TODO configurable size?
	stream := make(chan Event, 128)
	client := &Client{Stream: stream}

	go func() {
		defer resp.Body.Close()
		defer close(stream)

		client.Err = listenEvents(resp.Body, stream)
	}()

	return client, nil
}

func listenEvents(rc io.ReadCloser, stream chan Event) error {
	s := bufio.NewScanner(rc)
	s.Split(scanLines)

	var event Event

	for s.Scan() {
		line := s.Bytes()

		if len(line) == 0 {
			stream <- event
			event = Event{}
			continue
		}

		key, value := tokenize(line)

		switch string(key) {
		case "event":
			event.Name = string(value)
		case "data":
			event.Data = string(value)
		case "id":
			event.Id = string(value)
		case "retry":
			// TODO
		default:
			panic("Invalid token '" + key + "'")
		}
	}

	return s.Err()
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
