package dispatcher_client

import (
	"net"

	"github.com/xiaonanln/vacuum/netutil"
	"github.com/xiaonanln/vacuum/proto"
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
	req := proto.RegisterVacuumServerReq{}
	dc.send(proto.REGISTER_VACUUM_SERVER_REQ, &req)
}
