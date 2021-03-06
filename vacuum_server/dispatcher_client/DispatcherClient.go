package dispatcher_client

import (
	"net"

	"github.com/xiaonanln/vacuum/common"
	. "github.com/xiaonanln/vacuum/proto"
	"github.com/xiaonanln/vacuum/vlog"
)

type DispatcherRespHandler interface {
	HandleDispatcherResp_RegisterVacuumServer(serverIDs []int)
	HandleDispatcherResp_CreateString(name string, stringID string, args []interface{})
	HandleDispatcherResp_LoadString(name string, stringID string, args []interface{})
	HandleDispatcherResp_DeclareService(stringID string, serviceName string)
	HandleDispatcherResp_SendStringMessage(stringID string, msg common.StringMessage, tag uint)
	HandleDispatcherResp_CloseString(stringID string)
	HandleDispatcherResp_DelString(stringID string)
	HandleDispatcherResp_OnMigrateString(name string, stringID string, initArgs []interface{}, data map[string]interface{}, extraMigrateInfo map[string]interface{})
	HandleDispatcherResp_MigrateString(stringID string)
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

func (dc *DispatcherClient) SendStringMessage(stringID string, msg interface{}, tag uint) error {
	req := StringMessageRelay{
		Tag: tag,
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

func (dc *DispatcherClient) SendLoadStringReq(name string, stringID string, args []interface{}) error {
	req := LoadStringReq{
		Name:     name,
		StringID: stringID,
		Args:     args,
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

func (dc *DispatcherClient) SendMigrateStringReq(name string, stringID string, serverID int, towardsID string, initArgs []interface{}, data map[string]interface{}, extraMigrateInfo map[string]interface{}) error {
	req := MigrateStringReq{
		Name:             name,
		StringID:         stringID,
		ServerID:         serverID,
		TowardsID:        towardsID,
		Args:             initArgs,
		Data:             data,
		ExtraMigrateInfo: extraMigrateInfo,
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

func (dc *DispatcherClient) HandleMsg(msg *Message, pktSize uint32, msgtype MsgType_t) (err error) {
	vlog.Debug("<<< HandleMsg: size %v: %s", pktSize, MsgTypeToString(msgtype))
	payload := msg[PREPAYLOAD_SIZE:pktSize]
	if msgtype == START_MIGRATE_STRING_REQ {
		err = dc.handleStartMigrateStringReq(payload)
	} else if msgtype == MIGRATE_STRING_REQ {
		// migrate string to this server
		err = dc.handleMigrateStringReq(payload)
	} else if msgtype == CREATE_STRING_REQ {
		// create real string instance
		err = dc.handleCreateStringReq(payload)
	} else if msgtype == DECLARE_SERVICE_REQ {
		// declare service
		err = dc.handleDeclareServiceReq(payload)
	} else if msgtype == STRING_DEL_REQ {
		err = dc.handleStringDelReq(payload)
	} else if msgtype == LOAD_STRING_REQ {
		err = dc.handleLoadStringReq(payload)
	} else if msgtype == REGISTER_VACUUM_SERVER_RESP {
		err = dc.handleRegisterVacuumServerResp(payload)
	} else {
		vlog.Panicf("serveDispatcherClient: invalid msg type: %v", msgtype)
	}
	msg.Release()
	return
}

func (dc *DispatcherClient) HandleRelayMsg(msg *Message, pktSize uint32, targetID string) (err error) {
	var msgType MsgType_t = MsgType_t(NETWORK_ENDIAN.Uint16(msg[SIZE_FIELD_SIZE+STRING_ID_SIZE : SIZE_FIELD_SIZE+STRING_ID_SIZE+TYPE_FIELD_SIZE]))

	vlog.Debug("<<< HandleRelayMsg: size %v: %s", pktSize, MsgTypeToString(msgType))
	defer vlog.Debug("<<< HandleRelayMsg END")

	payload := msg[RELAY_PREPAYLOAD_SIZE:pktSize]
	if msgType == STRING_MESSAGE_RELAY {
		err = dc.handleSendStringRelay(targetID, payload)
	} else if msgType == CLOSE_STRING_RELAY {
		err = dc.handleCloseStringRelay(targetID)
	} else {
		vlog.Panicf("invalid msg type: %v", msgType)
	}
	msg.Release()
	return
}

func (dc *DispatcherClient) handleSendStringRelay(targetID string, payload []byte) error {
	var pkt StringMessageRelay
	err := MSG_PACKER.UnpackMsg(payload, &pkt)
	if err != nil {
		return err
	}

	dispatcherRespHandler.HandleDispatcherResp_SendStringMessage(targetID, pkt.Msg, pkt.Tag)
	return nil
}

func (dc *DispatcherClient) handleCloseStringRelay(targetID string) error {
	dispatcherRespHandler.HandleDispatcherResp_CloseString(targetID)
	return nil
}

func (dc *DispatcherClient) handleRegisterVacuumServerResp(payload []byte) error {
	var resp RegisterVacuumServerResp
	if err := MSG_PACKER.UnpackMsg(payload, &resp); err != nil {
		return err
	}

	dispatcherRespHandler.HandleDispatcherResp_RegisterVacuumServer(resp.ServerIDS)
	return nil
}

func (dc *DispatcherClient) handleCreateStringReq(payload []byte) error {
	var req CreateStringReq
	if err := MSG_PACKER.UnpackMsg(payload, &req); err != nil {
		return err
	}

	dispatcherRespHandler.HandleDispatcherResp_CreateString(req.Name, req.StringID, req.Args)
	return nil
}

func (dc *DispatcherClient) handleDeclareServiceReq(payload []byte) error {
	var req DeclareServiceReq
	err := MSG_PACKER.UnpackMsg(payload, &req)
	if err != nil {
		return err
	}

	dispatcherRespHandler.HandleDispatcherResp_DeclareService(req.StringID, req.ServiceName)
	return nil
}

func (dc *DispatcherClient) handleStringDelReq(payload []byte) error {
	var req StringDelReq
	if err := MSG_PACKER.UnpackMsg(payload, &req); err != nil {
		return err
	}

	dispatcherRespHandler.HandleDispatcherResp_DelString(req.StringID)
	return nil
}

func (dc *DispatcherClient) handleLoadStringReq(payload []byte) error {
	var req LoadStringReq
	if err := MSG_PACKER.UnpackMsg(payload, &req); err != nil {
		return err
	}

	dispatcherRespHandler.HandleDispatcherResp_LoadString(req.Name, req.StringID, req.Args)
	return nil
}

func (dc *DispatcherClient) handleStartMigrateStringReq(payload []byte) error {
	// Received start-migrate from dispatcher, now we start the real migrate progress
	var req StartMigrateStringReq
	if err := MSG_PACKER.UnpackMsg(payload, &req); err != nil {
		return err
	}

	dispatcherRespHandler.HandleDispatcherResp_MigrateString(req.StringID)
	return nil
}

func (dc *DispatcherClient) handleMigrateStringReq(payload []byte) error {
	var req MigrateStringReq
	if err := MSG_PACKER.UnpackMsg(payload, &req); err != nil {
		return err
	}

	dispatcherRespHandler.HandleDispatcherResp_OnMigrateString(req.Name, req.StringID, req.Args, req.Data, req.ExtraMigrateInfo)
	return nil
}
