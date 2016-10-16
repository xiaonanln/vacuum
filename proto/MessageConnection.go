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
	MAX_MESSAGE_SIZE        = 1 * 1024 * 1024
	MESSAGE_SIZE_FIELD_SIZE = 4
	MESSAGE_TYPE_FIELD_SIZE = 2
	MESSAGE_PREPAYLOAD_SIZE = MESSAGE_SIZE_FIELD_SIZE + MESSAGE_TYPE_FIELD_SIZE
)

var (
	NETWORK_ENDIAN = binary.LittleEndian
	msgPacker      = MessagePackMsgPacker{}
)

func init() {
	if MAX_MESSAGE_SIZE > msgbufpool.MSGBUF_SIZE {
		log.Panicf("MAX_MESSAGE_SIZE must be less than msgbufpool.MSGBUF_SIZE!")
	}
}

type MessageConnection struct {
	netutil.BinaryConnection
}

func NewMessageConnection(conn net.Conn) MessageConnection {
	return MessageConnection{BinaryConnection: netutil.NewBinaryConnection(conn)}
}

type MsgPacketInfo struct {
	Msgbuf  *msgbufpool.Msgbuf_t
	MsgType MsgType_t
	Payload []byte
}

func (mc MessageConnection) SendMsg(mt MsgType_t, msg interface{}) error {
	var msgbuf [MAX_MESSAGE_SIZE]byte

	NETWORK_ENDIAN.PutUint16((msgbuf)[MESSAGE_SIZE_FIELD_SIZE:MESSAGE_SIZE_FIELD_SIZE+MESSAGE_TYPE_FIELD_SIZE], uint16(mt))
	payloadBuf := (msgbuf)[MESSAGE_PREPAYLOAD_SIZE:MESSAGE_PREPAYLOAD_SIZE]
	payloadCap := cap(payloadBuf)
	payloadBuf, err := msgPacker.PackMsg(msg, payloadBuf)
	if err != nil {
		return err
	}

	payloadLen := len(payloadBuf)
	if payloadLen > payloadCap {
		// exceed payload
		return fmt.Errorf("MessageConnection: message paylaod too large(%d): %v", payloadLen, msg)
	}

	var pktSize uint32 = uint32(payloadLen + MESSAGE_PREPAYLOAD_SIZE)
	NETWORK_ENDIAN.PutUint32((msgbuf)[:MESSAGE_SIZE_FIELD_SIZE], pktSize)
	log.Printf("Send message: size=%v, type=%v: %v", pktSize, mt, msg)
	return mc.SendAll((msgbuf)[:pktSize])
}

func (mc MessageConnection) RecvMsgPacket(pinfo *MsgPacketInfo) error {
	var _sizeBuf [4]byte
	pktSizeBuf := _sizeBuf[:]
	err := mc.RecvAll(pktSizeBuf)
	if err != nil {
		return err
	}

	var pktSize uint32 = NETWORK_ENDIAN.Uint32(pktSizeBuf)
	if pktSize > MAX_MESSAGE_SIZE {
		// pkt size is too large
		return fmt.Errorf("message packet too large: %v", pktSize)
	}

	msgbuf := msgbufpool.GetMsgBuf()
	err = mc.RecvAll((*msgbuf)[MESSAGE_SIZE_FIELD_SIZE:pktSize])
	if err != nil {
		msgbufpool.PutMsgBuf(msgbuf) // put it back on error
		return nil
	}

	var msgtype MsgType_t
	msgtype = MsgType_t(NETWORK_ENDIAN.Uint16((*msgbuf)[MESSAGE_SIZE_FIELD_SIZE : MESSAGE_SIZE_FIELD_SIZE+MESSAGE_TYPE_FIELD_SIZE]))
	//msgPacker.UnpackMsg((*msgbuf)[MESSAGE_PREPAYLOAD_SIZE:pktSize),

	pinfo.Msgbuf = msgbuf
	pinfo.MsgType = msgtype
	pinfo.Payload = (*msgbuf)[MESSAGE_PREPAYLOAD_SIZE:pktSize]

	return nil
}
