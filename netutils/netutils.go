package netutils

import "net"

func IsTemporaryNetworkError(err error) bool {
	netErr, ok := err.(net.Error)
	if !ok {
		return false
	}
	return netErr.Temporary() || netErr.Timeout()
}
