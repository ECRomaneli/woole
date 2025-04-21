package sniffer

import (
	"net/http"
	web "woole/web/client"

	"github.com/ecromaneli-golang/http/webserver"
)

func setupServer() *webserver.Server {
	server := webserver.NewServerWithFS(http.FS(web.EmbeddedFS))
	server.Logger().SetLogLevelStr(config.SnifferLogLevel)

	server.FileServer()
	server.Get("/record/stream", connHandler)
	server.Get("/record/{id}/response/body", responseBodyHandler)
	server.Get("/record/{id}/replay", replayHandler)
	server.Post("/record/request", newRequestHandler)
	server.Delete("/record", clearHandler)

	return server
}
