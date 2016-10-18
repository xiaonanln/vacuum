package client_proxy

import (
	"net"

	"log"

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

func (cp *ClientProxy) Serve() {
	defer cp.Close()
	defer onClientProxyClose(cp)

	var err error

	log.Printf("New dispatcher client: %s", cp)
	var msgPacketInfo proto.MsgPacketInfo
	for {

		err = cp.RecvMsgPacket(&msgPacketInfo)
		if err != nil {
			// error
			break
		}

		log.Printf("dispatcher: received client msg: %v", msgPacketInfo)

		msgType := msgPacketInfo.MsgType
		if msgType == proto.SEND_STRING_MESSAGE_REQ {
			cp.handleSendStringMessageReq(msgPacketInfo.Payload)
		} else if msgType == proto.CREATE_STRING_REQ {
			cp.handleCreateStringReq(msgPacketInfo.Payload)
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

	log.Printf("%s.handleSendStringMessageReq %T %v, target serve r%", cp, req, req)
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

	log.Printf("%s.handleCreateStringReq %T %v, choose random server: %s", cp, req, req, chooseServer)
	resp := proto.CreateStringResp{
		Name:     req.Name,
		StringID: stringID,
	}

	chooseServer.SendMsg(proto.CREATE_STRING_RESP, &resp)
}

func (cp *ClientProxy) handleRegisterVacuumServerReq(data []byte) {
	var req proto.RegisterVacuumServerReq
	proto.MSG_PACKER.UnpackMsg(data, &req)
	log.Printf("%s.handleRegisterVacuumServerReq %T %v", cp, req, req)
	registerClientProxyInfo(cp, req.ServerID)
}

func (cp *ClientProxy) handleDeclareServiceReq(data []byte) {
	var req proto.DeclareServiceReq
	proto.MSG_PACKER.UnpackMsg(data, &req)
	log.Printf("%s.handleDeclareServiceReq %T %v", cp, req, req)

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
