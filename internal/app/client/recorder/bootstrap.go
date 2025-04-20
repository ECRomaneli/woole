package recorder

import (
	"os"
	"woole/internal/app/client/app"
	"woole/internal/app/client/recorder/adt"

	"github.com/ecromaneli-golang/console/logger"
)

const StatusInternalProxyError = 999

var config = app.ReadConfig()
var log = logger.New("recorder")

var records = adt.NewRecords(uint(config.MaxRecords))

func Start() {
	if config.IsStandalone {
		startStandalone()
	} else {
		startConnectionWithServer(onTunnelStart)
	}
	if app.HasSession() && !config.DisableSnifferOnlyMode {
		log.Warn("Tunnel connection closed, entering sniffer-only mode")
	} else {
		os.Exit(1)
	}
}
