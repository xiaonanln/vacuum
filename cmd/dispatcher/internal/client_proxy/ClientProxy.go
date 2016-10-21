package client_proxy

import (
	"net"

	log "github.com/Sirupsen/logrus"

	"fmt"

	"runtime/debug"

	"github.com/xiaonanln/vacuum/msgbufpool"
	"github.com/xiaonanln/vacuum/proto"
)

type ClientProxy struct {
	proto.MessageConnection
	ServerID int
}

func NewClientProxy(conn net.Conn) *ClientProxy {
	return &ClientProxy{
		MessageConnection: proto.NewMessageConnection(conn),
		ServerID:          0, // to be registered
	}
}

func (cp *ClientProxy) String() string {
	return fmt.Sprintf("Server<%d>", cp.ServerID)
}

func (cp *ClientProxy) Serve() {
	defer func() {
		cp.Close()
		onClientProxyClose(cp)

		err := recover()
		if err != nil {
			log.Errorf("Client %s paniced with error: %v", cp, err)
			debug.PrintStack()
		}
	}()

	var err error

	log.Infof("New dispatcher client: %s", cp)
	var msgPacketInfo proto.MsgPacketInfo
	for {

		err = cp.RecvMsgPacket(&msgPacketInfo)
		if err != nil {
			// error
			break
		}

		log.Debugf("dispatcher: received client msg: %v", msgPacketInfo)

		msgType := msgPacketInfo.MsgType
		if msgType == proto.SEND_STRING_MESSAGE_REQ {
			cp.handleSendStringMessageReq(msgPacketInfo.Payload)
		} else if msgType == proto.CREATE_STRING_REQ {
			cp.handleCreateStringReq(msgPacketInfo.Payload)
		} else if msgType == proto.CREATE_STRING_LOCALLY_REQ {
			cp.handleCreateStringLocallyReq(msgPacketInfo.Payload)
		} else if msgType == proto.REGISTER_VACUUM_SERVER_REQ {
			cp.handleRegisterVacuumServerReq(msgPacketInfo.Payload)
		} else if msgType == proto.DECLARE_SERVICE_REQ {
			cp.handleDeclareServiceReq(msgPacketInfo.Payload)
		} else {
			log.Panicf("ERROR: unknown dispatcher request type=%v", msgType)
		}

		msgbufpool.PutMsgBuf(msgPacketInfo.Msgbuf)
	}
}

func (cp *ClientProxy) handleSendStringMessageReq(data []byte) {
	var req proto.SendStringMessageReq
	proto.MSG_PACKER.UnpackMsg(data, &req)

	targetStringID := req.StringID
	resp := proto.SendStringMessageResp{
		StringID: targetStringID,
		Msg:      req.Msg,
	}

	serverID := getStringLocation(targetStringID)
	chooseServer := getClientProxy(serverID)

	log.Debugf("%s.handleSendStringMessageReq %T %v, target server %s", cp, req, req, chooseServer)
	chooseServer.SendMsg(proto.SEND_STRING_MESSAGE_RESP, &resp)
}

func (cp *ClientProxy) handleCreateStringReq(data []byte) {
	var req proto.CreateStringReq
	proto.MSG_PACKER.UnpackMsg(data, &req)

	// choose one server for create string

	chooseServer := getRandomClientProxy()
	// save the stringID with the serverID
	stringID := req.StringID
	setStringLocation(stringID, chooseServer.ServerID)

	log.Debugf("%s.handleCreateStringReq %T %v, choose random server: %s", cp, req, req, chooseServer)
	resp := proto.CreateStringResp{
		Name:     req.Name,
		StringID: stringID,
	}

	chooseServer.SendMsg(proto.CREATE_STRING_RESP, &resp)
}

func (cp *ClientProxy) handleCreateStringLocallyReq(data []byte) {
	var req proto.CreateStringLocallyReq
	proto.MSG_PACKER.UnpackMsg(data, &req)

	// choose one server for create string

	stringID := req.StringID
	setStringLocation(stringID, cp.ServerID)
	log.Debugf("%s.handleCreateStringLocallyReq %T %v", cp, req, req)
}

func (cp *ClientProxy) handleRegisterVacuumServerReq(data []byte) {
	var req proto.RegisterVacuumServerReq
	proto.MSG_PACKER.UnpackMsg(data, &req)
	log.Debugf("%s.handleRegisterVacuumServerReq %T %v", cp, req, req)
	registerClientProxyInfo(cp, req.ServerID)
}

func (cp *ClientProxy) handleDeclareServiceReq(data []byte) {
	var req proto.DeclareServiceReq
	proto.MSG_PACKER.UnpackMsg(data, &req)
	log.Debugf("%s.handleDeclareServiceReq %T %v", cp, req, req)

	// the the declare of service to all clients
	resp := proto.DeclareServiceResp{
		StringID:    req.StringID,
		ServiceName: req.ServiceName,
	}
	clientProxiesLock.RLock()
	for _, clientProxy := range clientProxes {
		clientProxy.SendMsg(proto.DECLARE_SERVICE_RESP, &resp)
	}
	clientProxiesLock.RUnlock()
}
