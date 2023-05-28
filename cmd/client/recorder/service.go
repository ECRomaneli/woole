package recorder

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"woole/cmd/client/app"
	"woole/cmd/client/recorder/adt"
	pb "woole/shared/payload"
	"woole/shared/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Replay(request *pb.Request) {
	record := adt.NewRecord(request, adt.REPLAY)
	record.Response = proxyRequest(record.Request)
	records.AddRecordAndCallListeners(record)

	if log.IsInfoEnabled() {
		log.Info(record.ToString(26))
	}
}

func GetRecords() *adt.Records {
	return records
}

func onTunnelStart(client pb.TunnelClient, ctx context.Context, cancelCtx context.CancelFunc) error {
	defer cancelCtx()

	// Start the tunnel stream
	stream, err := client.Tunnel(ctx)
	if !handleGRPCErrors(err) {
		return err
	}

	// Send the handshake
	stream.Send(&pb.ClientMessage{Handshake: config.GetHandshake()})

	// Receive the session
	serverMsg, err := stream.Recv()
	if !handleGRPCErrors(err) {
		return err
	}

	app.SetSession(serverMsg.Session)

	// Listen for requests and send responses asynchronously
	for {
		serverMsg, err := stream.Recv()

		if err != nil {
			if !handleGRPCErrors(err) {
				return err
			}
			continue
		}

		go handleServerRecord(stream, serverMsg.Record)
	}
}

func handleServerRecord(stream pb.Tunnel_TunnelClient, record *pb.Record) {
	defer catchAllErrors(record)

	switch record.Step {
	case pb.Step_SEND_REQUEST:
		handleServerRequest(stream, record)
	case pb.Step_SEND_SERVER_RESPONSE:
		handleServerResponse(stream, record)
	default:
		log.Error("Record Step Not Allowed")
	}
}

func handleServerRequest(stream pb.Tunnel_TunnelClient, serverRecord *pb.Record) {
	record := adt.EnhanceRecord(serverRecord)

	doRequest(record)
	records.AddRecordAndCallListeners(record)
	err := stream.Send(&pb.ClientMessage{Record: record.Record})

	if log.IsInfoEnabled() {
		log.Info(record.ToString(26))
	}

	if !handleGRPCErrors(err) {
		log.Error("Failed to send response for Record[", record.Id, "].", err)
	}
}

func handleServerResponse(stream pb.Tunnel_TunnelClient, serverRecord *pb.Record) {
	record := records.FindById(serverRecord.Id)
	log.Debug("Received a STEP.RECEIVE_SERVER_RESPONSE for record [", record, "]")
}

func doRequest(record *adt.Record) {
	record.Step = pb.Step_RECEIVE_RESPONSE
	replaceUrlHeaderByCustomUrl(record.Request.Header, "Origin")
	replaceUrlHeaderByCustomUrl(record.Request.Header, "Referer")

	handleWooleParamIfExists(record)
	if record.Response != nil {
		return
	}

	record.Response = proxyRequest(record.Request)
	handleRedirections(record)
}

func proxyRequest(req *pb.Request) *pb.Response {
	// Redirect and record the response
	recorder := httptest.NewRecorder()
	elapsed := util.Timer(func() {
		proxyHandler.ServeHTTP(recorder, req.ToHTTPRequest())
	})

	// Save req and res data
	return (&pb.Response{}).FromResponseRecorder(recorder, elapsed)
}

func handleRedirections(record *adt.Record) {
	location := record.Response.GetHttpHeader().Get("location")

	if location == "" {
		return
	}

	record.OriginalResponse = record.Response
	record.Response = record.OriginalResponse.Clone()

	urlParam := &adt.CustomPathParam{Redirect: adt.Redirect{
		RecordId: record.Id,
		Action:   adt.CONTINUE,
	}}

	params := make(map[string]string)
	params["redirectUrl"] = location
	params["hostname"] = app.GetSessionWhenAvailable().Hostname
	params["originalUrl"] = urlParam.Serialize()

	if record.Type == adt.REDIRECT {
		params["enableCustomUrl"] = "false"
		params["customUrl"] = "#"
	} else {
		urlParam.Redirect.Action = adt.CHANGE_URL_HOST
		params["enableCustomUrl"] = "true"
		params["customUrl"] = urlParam.Serialize()
	}

	record.Response.Body = []byte(app.RedirectTemplate.Apply(params))
	record.Response.Code = http.StatusOK

	httpHeader := record.Response.GetHttpHeader()
	httpHeader.Set("Content-Type", "text/html")
	httpHeader.Del("location")
	httpHeader.Set("Content-Length", strconv.Itoa(len(record.Response.Body)))
	record.Response.SetHttpHeader(httpHeader)
}

func handleWooleParamIfExists(record *adt.Record) {
	param, ok := adt.DeserializeCustomPathParam(record.Request.Url)
	if !ok {
		return
	}

	redirectRecord := records.FindById(param.Redirect.RecordId)
	if redirectRecord == nil {
		panic("Trying to access an invalid record [" + param.Redirect.RecordId + "]")
	}

	record.Type = adt.REDIRECT

	if param.Redirect.Action == adt.CONTINUE {
		record.Request = redirectRecord.Request.Clone()
		record.Response = redirectRecord.OriginalResponse
		return
	}

	location := redirectRecord.OriginalResponse.GetHttpHeader().Get("location")

	if len(location) == 0 {
		panic("There is no redirect URL")
	}

	newUrl, ok := util.ReplaceHostByUsingExampleStr(location, record.Request.Url)
	if !ok {
		panic("Error when trying to replace the host of [" + record.Request.Url + "]")
	}

	record.Request.Url = newUrl.String()
	record.Request.Path = newUrl.Path
}

func replaceUrlHeaderByCustomUrl(header map[string]*pb.StringList, headerName string) {
	if header == nil || header[headerName] == nil {
		return
	}

	rawUrl := header[headerName].Val[0]
	newUrl, ok := util.ReplaceHostByUsingExampleUrl(rawUrl, config.CustomUrl)

	if !ok {
		panic("Error when trying to replace the host of [" + rawUrl + "]")
	}

	header[headerName] = &pb.StringList{Val: []string{newUrl.String()}}
}

func handleGRPCErrors(err error) bool {
	if err == nil {
		return true
	}

	switch status.Code(err) {
	case codes.ResourceExhausted:
		log.Warn("Request discarded. Reason: Max size exceeded")
		return true
	default:
		return false
	}
}

func catchAllErrors(record *pb.Record) {
	err := recover()

	if err == nil {
		return
	}

	log.Error(err)
	// TODO: Improve error handling
}
