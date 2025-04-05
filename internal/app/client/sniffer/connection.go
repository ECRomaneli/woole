package sniffer

import (
	"net/http"
	"woole/web"

	"github.com/ecromaneli-golang/http/webserver"
)

func setupServer() *webserver.Server {
	server := webserver.NewServerWithFS(http.FS(web.EmbeddedFS))

	server.FileServer("/")
	server.Get("/record/stream", connHandler)
	server.Get("/record/{id}/response/body", responseBodyHandler)
	server.Get("/record/{id}/replay", replayHandler)
	server.Post("/record/request", newRequestHandler)
	server.Delete("/record", clearHandler)

	return server
}
