package constants

import "strconv"

const (
	DefaultTunnelPort = 8001
)

var (
	DefaultTunnelPortStr = strconv.FormatInt(int64(DefaultTunnelPort), 10)
)
