package dispatcher_client

import (
	"net"

	"github.com/xiaonanln/vacuum/common"
	. "github.com/xiaonanln/vacuum/proto"
	"github.com/xiaonanln/vacuum/vlog"
)

type DispatcherRespHandler interface {
	HandleDispatcherResp_CreateString(name string, stringID string, args []interface{})
	HandleDispatcherResp_LoadString(name string, stringID string)
	HandleDispatcherResp_DeclareService(stringID string, serviceName string)
	HandleDispatcherResp_SendStringMessage(stringID string, msg common.StringMessage)
	HandleDispatcherResp_CloseString(stringID string)
	HandleDispatcherResp_DelString(stringID string)
	HandleDispatcherResp_OnMigrateString(name string, stringID string, data map[string]interface{})
	HandleDispatcherResp_StartMigrateString(stringID string)
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

func (dc *DispatcherClient) SendLoadStringReq(name string, stringID string) error {
	req := LoadStringReq{
		Name:     name,
		StringID: stringID,
	}
	return dc.SendMsg(LOAD_STRING_REQ, &req)
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

func (dc *DispatcherClient) SendStartMigrateStringReq(stringID string) error {
	req := StartMigrateStringReq{
		StringID: stringID,
	}
	return dc.SendMsg(START_MIGRATE_STRING_REQ, &req)
}

func (dc *DispatcherClient) SendMigrateStringReq(name string, stringID string, serverID int, data map[string]interface{}) error {
	req := MigrateStringReq{
		Name:     name,
		StringID: stringID,
		ServerID: serverID,
		Data:     data,
	}

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
	vlog.Debug("<<< HandleMsg: size %v, msgtype %v", pktSize, msgtype)
	payload := msg[PREPAYLOAD_SIZE:pktSize]
	if msgtype == START_MIGRATE_STRING_RESP {
		return dc.handleStartMigrateStringResp(payload)
	} else if msgtype == MIGRATE_STRING_RESP {
		// migrate string to this server
		return dc.handleMigrateStringResp(payload)
	} else if msgtype == CREATE_STRING_RESP {
		// create real string instance
		return dc.handleCreateStringResp(payload)
	} else if msgtype == DECLARE_SERVICE_RESP {
		// declare service
		return dc.handleDeclareServiceResp(payload)
	} else if msgtype == STRING_DEL_RESP {
		return dc.handleStringDelResp(payload)
	} else if msgtype == LOAD_STRING_RESP {
		return dc.handleLoadStringResp(payload)
	} else {
		vlog.Panicf("serveDispatcherClient: invalid msg type: %v", msgtype)
		return nil
	}
}

func (dc *DispatcherClient) HandleRelayMsg(msg *Message, pktSize uint32, targetID string) error {
	var msgType MsgType_t = MsgType_t(NETWORK_ENDIAN.Uint16(msg[SIZE_FIELD_SIZE+STRING_ID_SIZE : SIZE_FIELD_SIZE+STRING_ID_SIZE+TYPE_FIELD_SIZE]))

	vlog.Debug("<<< HandleRelayMsg: size %v, msgtype %v", pktSize, msgType)
	defer vlog.Debug("<<< HandleRelayMsg END")

	payload := msg[RELAY_PREPAYLOAD_SIZE:pktSize]
	if msgType == STRING_MESSAGE_RELAY {
		return dc.handleSendStringRelay(targetID, payload)
	} else if msgType == CLOSE_STRING_RELAY {
		return dc.handleCloseStringRelay(targetID)
	} else {
		vlog.Panicf("invalid msg type: %v", msgType)
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

func (dc *DispatcherClient) handleLoadStringResp(payload []byte) error {
	var resp LoadStringResp
	if err := MSG_PACKER.UnpackMsg(payload, &resp); err != nil {
		return err
	}

	dispatcherRespHandler.HandleDispatcherResp_LoadString(resp.Name, resp.StringID)
	return nil
}

func (dc *DispatcherClient) handleStartMigrateStringResp(payload []byte) error {
	// Received start-migrate from dispatcher, now we start the real migrate progress
	var resp StartMigrateStringResp
	if err := MSG_PACKER.UnpackMsg(payload, &resp); err != nil {
		return err
	}

	dispatcherRespHandler.HandleDispatcherResp_StartMigrateString(resp.StringID)
	return nil
}

func (dc *DispatcherClient) handleMigrateStringResp(payload []byte) error {
	var resp MigrateStringResp
	if err := MSG_PACKER.UnpackMsg(payload, &resp); err != nil {
		return err
	}

	dispatcherRespHandler.HandleDispatcherResp_OnMigrateString(resp.Name, resp.StringID, resp.Data)
	return nil
}
