package constants

import "strconv"

const (
	DefaultTunnelPort = 9653
)

var (
	DefaultTunnelPortStr = strconv.FormatInt(int64(DefaultTunnelPort), 10)
)
