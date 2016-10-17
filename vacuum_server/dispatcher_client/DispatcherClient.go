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

func (dc *DispatcherClient) RegisterVacuumServer(serverID int) error {
	req := proto.RegisterVacuumServerReq{
		ServerID: serverID,
	}
	return dc.SendMsg(proto.REGISTER_VACUUM_SERVER_REQ, &req)
}

func (dc *DispatcherClient) SendStringMessage(sid string, msg interface{}) error {
	req := proto.SendStringMessageReq{
		SID: sid,
		Msg: msg,
	}
	return dc.SendMsg(proto.SEND_STRING_MESSAGE_REQ, &req)
}

func (dc *DispatcherClient) CreateString(name string) error {
	req := proto.CreateStringReq{
		Name: name,
	}
	return dc.SendMsg(proto.CREATE_STRING_REQ, &req)
}
