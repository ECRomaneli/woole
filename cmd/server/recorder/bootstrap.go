package recorder

import (
	"woole/cmd/server/app"
	"woole/cmd/server/recorder/adt"
	pb "woole/shared/payload"

	"github.com/ecromaneli-golang/console/logger"
)

var (
	config        = app.ReadConfig()
	log           = logger.New("recorder")
	clientManager = adt.NewClientManager()
)

type Tunnel struct{ pb.UnimplementedTunnelServer }

func Start() {
	serveTunnel()
	serveWebServer()
}
