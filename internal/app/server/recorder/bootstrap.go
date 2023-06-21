package recorder

import (
	"woole/internal/app/server/app"
	"woole/internal/app/server/recorder/adt"
	pb "woole/internal/pkg/payload"

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
