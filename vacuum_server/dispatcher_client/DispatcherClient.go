package dispatcher_client

import (
	"log"
	"net"

	"github.com/xiaonanln/vacuum/common"
	. "github.com/xiaonanln/vacuum/proto"
)

type DispatcherRespHandler interface {
	HandleDispatcherResp_CreateString(name string, stringID string, args []interface{})
	HandleDispatcherResp_DeclareService(stringID string, serviceName string)
	HandleDispatcherResp_SendStringMessage(stringID string, msg common.StringMessage)
	HandleDispatcherResp_CloseString(stringID string)
	HandleDispatcherResp_DelString(stringID string)
}

type DispatcherClient struct {
	MessageConnection
}

var (
	dispatcherRespHandler DispatcherRespHandler
)

func newDispatcherClient(conn net.Conn) *DispatcherClient {
	return &DispatcherClient{
		MessageConnection: NewMessageConnection(conn),
	}
}

func (dc *DispatcherClient) RegisterVacuumServer(serverID int) error {
	req := RegisterVacuumServerReq{
		ServerID: serverID,
	}
	return dc.SendMsg(REGISTER_VACUUM_SERVER_REQ, &req)
}

func (dc *DispatcherClient) SendStringMessage(stringID string, msg interface{}) error {
	req := StringMessageRelay{
		Msg: msg,
	}
	return dc.SendRelayMsg(stringID, STRING_MESSAGE_RELAY, &req)
}

func (dc *DispatcherClient) SendCreateStringReq(name string, stringID string, args []interface{}) error {
	req := CreateStringReq{
		Name:     name,
		StringID: stringID,
		Args:     args,
	}
	return dc.SendMsg(CREATE_STRING_REQ, &req)
}

func (dc *DispatcherClient) SendCreateStringLocallyReq(name string, stringID string) error {
	req := CreateStringLocallyReq{
		Name:     name,
		StringID: stringID,
	}
	return dc.SendMsg(CREATE_STRING_LOCALLY_REQ, &req)
}

func (dc *DispatcherClient) SendDeclareServiceReq(stringID string, serviceName string) error {
	req := DeclareServiceReq{
		StringID:    stringID,
		ServiceName: serviceName,
	}
	return dc.SendMsg(DECLARE_SERVICE_REQ, &req)
}

func (dc *DispatcherClient) SendStringDelReq(stringID string) error {
	req := StringDelReq{
		StringID: stringID,
	}
	return dc.SendMsg(STRING_DEL_REQ, &req)
}

func (dc *DispatcherClient) SendMigrateStringReq(stringID string, serverID int) error {
	req := MigrateStringReq{}
	return dc.SendMsg(MIGRATE_STRING_REQ, &req)
}

func (dc *DispatcherClient) RelayCloseString(stringID string) error {
	req := CloseStringRelay{}
	return dc.SendRelayMsg(stringID, CLOSE_STRING_RELAY, &req)
}

//
//func (mc MessageConnection) SendRawRelayMsg(targetID string, msgTypeAndPayload []byte) error {
//	var pkgSize uint32 = SIZE_FIELD_SIZE + STRING_ID_SIZE + len(msgTypeAndPayload)
//	var pkgSizeBuf [SIZE_FIELD_SIZE]byte
//	NETWORK_ENDIAN.PutUint32(pkgSizeBuf, pkgSize)
//
//}

func (dc *DispatcherClient) HandleMsg(msg *Message, pktSize uint32, msgtype MsgType_t) error {
	payload := msg[PREPAYLOAD_SIZE:pktSize]
	if msgtype == CREATE_STRING_RESP {
		// create real string instance
		return dc.handleCreateStringResp(payload)
	} else if msgtype == DECLARE_SERVICE_RESP {
		// declare service
		return dc.handleDeclareServiceResp(payload)
	} else if msgtype == STRING_DEL_RESP {
		return dc.handleStringDelResp(payload)
	} else {
		log.Panicf("serveDispatcherClient: invalid msg type: %v", msgtype)
		return nil
	}
}

func (dc *DispatcherClient) HandleRelayMsg(msg *Message, pktSize uint32, targetID string) error {
	var msgType MsgType_t = MsgType_t(NETWORK_ENDIAN.Uint16(msg[SIZE_FIELD_SIZE+STRING_ID_SIZE : SIZE_FIELD_SIZE+STRING_ID_SIZE+TYPE_FIELD_SIZE]))
	payload := msg[RELAY_PREPAYLOAD_SIZE:pktSize]
	if msgType == STRING_MESSAGE_RELAY {
		return dc.handleSendStringRelay(targetID, payload)
	} else if msgType == CLOSE_STRING_RELAY {
		return dc.handleCloseStringRelay(targetID)
	} else {
		log.Panicf("invalid msg type: %v", msgType)
		return nil
	}
}

func (dc *DispatcherClient) handleSendStringRelay(targetID string, payload []byte) error {
	var pkt StringMessageRelay
	err := MSG_PACKER.UnpackMsg(payload, &pkt)
	if err != nil {
		return err
	}

	dispatcherRespHandler.HandleDispatcherResp_SendStringMessage(targetID, pkt.Msg)
	return nil
}

func (dc *DispatcherClient) handleCloseStringRelay(targetID string) error {
	dispatcherRespHandler.HandleDispatcherResp_CloseString(targetID)
	return nil
}

func (dc *DispatcherClient) handleCreateStringResp(payload []byte) error {
	var resp CreateStringResp
	err := MSG_PACKER.UnpackMsg(payload, &resp)
	if err != nil {
		return err
	}

	dispatcherRespHandler.HandleDispatcherResp_CreateString(resp.Name, resp.StringID, resp.Args)
	return nil
}

func (dc *DispatcherClient) handleDeclareServiceResp(payload []byte) error {
	var resp DeclareServiceResp
	err := MSG_PACKER.UnpackMsg(payload, &resp)
	if err != nil {
		return err
	}

	dispatcherRespHandler.HandleDispatcherResp_DeclareService(resp.StringID, resp.ServiceName)
	return nil
}

func (dc *DispatcherClient) handleStringDelResp(payload []byte) error {
	var resp StringDelResp
	if err := MSG_PACKER.UnpackMsg(payload, &resp); err != nil {
		return err
	}

	dispatcherRespHandler.HandleDispatcherResp_DelString(resp.StringID)
	return nil
}
