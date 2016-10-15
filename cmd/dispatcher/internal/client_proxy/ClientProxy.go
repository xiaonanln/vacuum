package client_proxy

import (
	"net"

	"log"

	"github.com/xiaonanln/vacuum/proto"
)

type ClientProxy struct {
	proto.MessageConnection
}

func NewClientProxy(conn net.Conn) *ClientProxy {
	return &ClientProxy{
		MessageConnection: proto.NewMessageConnection(conn),
	}
}

func (cp *ClientProxy) Serve() {
	defer cp.Close()

	log.Printf("New dispatcher client: %s", cp)
}
