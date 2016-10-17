package dispatcher_client

import (
	"net"

	"github.com/xiaonanln/vacuum/proto"
)

type DispatcherRespHandler interface {
	HandleDispatcherResp_CreateString(name string)
	HandleDispatcherResp_DeclareService(sid string, serviceName string)
}

type DispatcherClient struct {
	proto.MessageConnection
}

var (
	dispatcherRespHandler DispatcherRespHandler
)

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

func (dc *DispatcherClient) SendCreateStringReq(name string) error {
	req := proto.CreateStringReq{
		Name: name,
	}
	return dc.SendMsg(proto.CREATE_STRING_REQ, &req)
}

func (dc *DispatcherClient) SendDeclareServiceReq(sid string, serviceName string) error {
	req := proto.DeclareServiceReq{
		StringID:    sid,
		ServiceName: serviceName,
	}
	return dc.SendMsg(proto.DECLARE_SERVICE_REQ, &req)
}
