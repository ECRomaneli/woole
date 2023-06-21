package recorder

import (
	"math"
	"woole/internal/app/client/app"
	"woole/internal/app/client/recorder/adt"
	pb "woole/internal/pkg/payload"

	"github.com/ecromaneli-golang/http/webserver"
)

func startStandalone() {
	app.SetSession(&pb.Session{
		ClientId:        "standalone",
		HttpPort:        config.HttpUrl.Port(),
		Hostname:        "localhost",
		MaxRequestSize:  math.MaxInt32,
		MaxResponseSize: math.MaxInt32,
	})

	server := webserver.NewServer()
	server.All(config.HttpUrl.Hostname()+"/**", recorderHandler)
	panic(server.ListenAndServe(":" + config.HttpUrl.Port()))
}

// REST = [ALL] /**
func recorderHandler(req *webserver.Request, res *webserver.Response) {
	record := adt.NewRecord((&pb.Request{}).FromHTTPRequest(req), adt.DEFAULT)
	doRequest(record)

	res.Headers(record.Response.GetHttpHeader()).Status(int(record.Response.Code)).Write(record.Response.Body)

	records.AddRecordAndPublish(record)

	if log.IsInfoEnabled() {
		log.Info(record.ToString(26))
	}
}
