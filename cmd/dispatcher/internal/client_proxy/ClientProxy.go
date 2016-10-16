package client_proxy

import (
	"net"

	"log"

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
	}
}
