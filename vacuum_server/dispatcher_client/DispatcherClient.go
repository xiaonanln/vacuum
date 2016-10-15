package dispatcher_client

import (
	"net"

	"github.com/xiaonanln/vacuum/netutil"
)

type DispatcherClient struct {
	netutil.BinaryConnection
}

func newDispatcherClient(conn net.Conn) *DispatcherClient {
	return &DispatcherClient{
		BinaryConnection: netutil.NewBinaryConnection(conn),
	}
}

func (dc *DispatcherClient) RegisterVacuumServer() {

}
