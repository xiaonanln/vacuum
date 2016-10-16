package dispatcher_client

import (
	"net"

	"github.com/xiaonanln/vacuum/proto"
)

type DispatcherClient struct {
	proto.MessageConnection
}

func newDispatcherClient(conn net.Conn) *DispatcherClient {
	return &DispatcherClient{
		MessageConnection: proto.NewMessageConnection(conn),
	}
}

func (dc *DispatcherClient) RegisterVacuumServer() {
	req := proto.RegisterVacuumServerReq{}
	dc.SendMsg(proto.REGISTER_VACUUM_SERVER_REQ, &req)
}

func (dc *DispatcherClient) SendStringMessage(sid string, msg interface{}) {
	req := proto.SendStringMessageReq{
		SID: sid,
		Msg: msg,
	}
	dc.SendMsg(proto.SEND_STRING_MESSAGE_REQ, &req)
}
