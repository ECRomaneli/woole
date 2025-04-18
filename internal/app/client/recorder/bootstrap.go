package recorder

import (
	"woole/internal/app/client/app"
	"woole/internal/app/client/recorder/adt"

	"github.com/ecromaneli-golang/console/logger"
)

const StatusInternalProxyError = 999

var config = app.ReadConfig()
var log = logger.New("recorder")

var records = adt.NewRecords(uint(config.MaxRecords))
var proxyHandler = CreateProxyHandler()

func Start() {
	if config.IsStandalone {
		startStandalone()
	} else {
		startConnectionWithServer(onTunnelStart)
	}
}
