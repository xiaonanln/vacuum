package proto

import (
	"net"

	"encoding/binary"

	"fmt"

	"log"

	"github.com/xiaonanln/vacuum/msgbufpool"
	"github.com/xiaonanln/vacuum/netutil"
)

const (
	MESSAGE_SIZE_FIELD_SIZE = 4
	MESSAGE_TYPE_FIELD_SIZE = 2
)

var (
	msgPacker = MessagePackMsgPacker{}
)

type MessageConnection struct {
	netutil.BinaryConnection
}

func NewMessageConnection(conn net.Conn) MessageConnection {
	return MessageConnection{BinaryConnection: netutil.NewBinaryConnection(conn)}
}

func (mc MessageConnection) SendMsg(mt MsgType_t, msg interface{}) error {
	msgbuf := msgbufpool.GetMsgBuf()
	binary.LittleEndian.PutUint16((*msgbuf)[MESSAGE_SIZE_FIELD_SIZE:MESSAGE_SIZE_FIELD_SIZE+MESSAGE_TYPE_FIELD_SIZE], uint16(mt))
	payloadBuf := (*msgbuf)[MESSAGE_SIZE_FIELD_SIZE+MESSAGE_TYPE_FIELD_SIZE : MESSAGE_SIZE_FIELD_SIZE+MESSAGE_TYPE_FIELD_SIZE]
	payloadCap := cap(payloadBuf)
	payloadBuf, err := msgPacker.PackMsg(msg, payloadBuf)
	if err != nil {
		msgbufpool.PutMsgBuf(msgbuf) // put msgbuf back
		return err
	}

	payloadLen := len(payloadBuf)
	if payloadLen > payloadCap {
		// exceed payload
		msgbufpool.PutMsgBuf(msgbuf) // put msgbuf back
		return fmt.Errorf("MessageConnection: message paylaod too large(%d): %v", payloadLen, msg)
	}

	var size uint32 = uint32(payloadLen + MESSAGE_TYPE_FIELD_SIZE)
	binary.LittleEndian.PutUint32((*msgbuf)[:MESSAGE_SIZE_FIELD_SIZE], size)
	log.Printf("Send message: size=%v, type=%v: %v", size, mt, msg)
	return mc.SendAll((*msgbuf)[:MESSAGE_SIZE_FIELD_SIZE+MESSAGE_TYPE_FIELD_SIZE+payloadLen])
}
