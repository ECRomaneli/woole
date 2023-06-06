package dashboard

import (
	"io/fs"
	"net/http"
	"woole/cmd/client/app"

	"github.com/ecromaneli-golang/http/webserver"
)

func setupServer() *webserver.Server {
	staticFolder, _ := fs.Sub(app.EmbeddedFS, "static")
	server := webserver.NewServerWithFS(http.FS(staticFolder))

	server.FileServer("/")
	server.Get("/record/stream", connHandler)
	server.Get("/record/{id}/response/body", responseBodyHandler)
	server.Get("/record/{id}/replay", replayHandler)
	server.Get("/record/{id}/request/curl", curlHandler)
	server.Post("/record/request", newRequestHandler)
	server.Delete("/record", clearHandler)

	return server
}
