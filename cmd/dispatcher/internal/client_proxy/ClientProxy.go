package client_proxy

import (
	"net"

	"log"

	"github.com/xiaonanln/vacuum/netutil"
)

type ClientProxy struct {
	netutil.BinaryConnection
}

func NewClientProxy(conn net.Conn) *ClientProxy {
	return &ClientProxy{
		BinaryConnection: netutil.NewBinaryConnection(conn),
	}
}

func (cp *ClientProxy) Serve() {
	defer cp.Close()

	log.Printf("New dispatcher client: %s", cp)
}
