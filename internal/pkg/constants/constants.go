package constants

import "strconv"

const (
	DefaultTunnelPort        = 9653
	ForwardedToHeader        = "X-Woole-Forwarded-To"
	DefaultProxyPort         = "80"
	DefaultSnifferPort       = "8000"
	DefaultStandalonePort    = "8080"
	DefaultStandaloneMessage = "[<hostname>]:<port>"
	DefaultCustomUrlMessage  = "[<scheme>://]<hostname>[:<port>]"
	ClientToken              = "{client}"
)

var (
	DefaultTunnelPortStr = strconv.FormatInt(int64(DefaultTunnelPort), 10)
)
