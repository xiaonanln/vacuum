package gameserver

import (
	"net"

	"github.com/xiaonanln/vacuum/netutil"
)

type GSClient struct {
	netutil.BinaryConnection
	conn net.Conn
}

func newGSClient(conn net.Conn) *GSClient {
	return &GSClient{
		conn: conn,
	}
}

func (client *GSClient) serve() {

}
