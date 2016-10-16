package client_proxy

import (
	"net"

	"log"

	"github.com/xiaonanln/vacuum/msgbufpool"
	"github.com/xiaonanln/vacuum/proto"
)

var (
	msgPacker = proto.MessagePackMsgPacker{}
)

type ClientProxy struct {
	proto.MessageConnection
}

func NewClientProxy(conn net.Conn) *ClientProxy {
	return &ClientProxy{
		MessageConnection: proto.NewMessageConnection(conn),
	}
}

func (cp *ClientProxy) Serve() {
	defer cp.Close()

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
		} else {
			log.Printf("ERROR: unknown dispatcher request type=%v", msgType)
		}

		msgbufpool.PutMsgBuf(msgPacketInfo.Msgbuf)
	}

}

func (cp *ClientProxy) handleSendStringMessageReq(data []byte) {
	var req proto.SendStringMessageReq
	msgPacker.UnpackMsg(data, &req)
	log.Printf("%s.handleSendStringMessageReq %T %v", cp, req, req)
}

func (cp *ClientProxy) handleCreateStringReq(data []byte) {
	var req proto.CreateStringReq
	msgPacker.UnpackMsg(data, &req)
	log.Printf("%s.handleCreateStringReq %T %v", cp, req, req)

	// choose one server for create string

}

func (cp *ClientProxy) handleRegisterVacuumServerReq(data []byte) {
	var req proto.RegisterVacuumServerReq
	msgPacker.UnpackMsg(data, &req)
	log.Println("%s.handleRegisterVacuumServerReq %T %v", cp, req, req)
}
