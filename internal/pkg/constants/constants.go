package constants

import "strconv"

const (
	DefaultTunnelPort = 9653
	ForwardedToHeader = "X-Woole-Forwarded-To"
)

var (
	DefaultTunnelPortStr = strconv.FormatInt(int64(DefaultTunnelPort), 10)
)
