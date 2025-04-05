package recorder

import (
	"woole/internal/app/server/app"
	"woole/internal/app/server/recorder/adt"
	"woole/internal/pkg/tunnel"

	"github.com/ecromaneli-golang/console/logger"
)

var (
	config        = app.ReadConfig()
	log           = logger.New("recorder")
	clientManager = adt.NewClientManager()
)

type Tunnel struct {
	tunnel.UnimplementedTunnelServer
}

func Start() {
	serveTunnel()
	serveWebServer()
}
