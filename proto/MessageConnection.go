package proto

import (
	"net"

	"encoding/binary"

	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/xiaonanln/vacuum/msgbufpool"
	"github.com/xiaonanln/vacuum/netutil"
	"github.com/xiaonanln/vacuum/uuid"
)

const (
	MAX_MESSAGE_SIZE = 1 * 1024 * 1024
	SIZE_FIELD_SIZE  = 4
	TYPE_FIELD_SIZE  = 2
	PREPAYLOAD_SIZE  = SIZE_FIELD_SIZE + TYPE_FIELD_SIZE

	STRING_ID_SIZE        = uuid.UUID_LENGTH
	RELAY_PREPAYLOAD_SIZE = SIZE_FIELD_SIZE + STRING_ID_SIZE + TYPE_FIELD_SIZE
)

var (
	NETWORK_ENDIAN = binary.LittleEndian
)

func init() {
	if MAX_MESSAGE_SIZE > msgbufpool.MSGBUF_SIZE {
		log.Panicln("MAX_MESSAGE_SIZE must be less than msgbufpool.MSGBUF_SIZE!")
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

// Send msg to/from dispatcher
// Message format: [size*4B][type*2B][payload*NB]
func (mc MessageConnection) SendMsg(mt MsgType_t, msg interface{}) error {
	var msgbuf [MAX_MESSAGE_SIZE]byte

	NETWORK_ENDIAN.PutUint16((msgbuf)[SIZE_FIELD_SIZE:SIZE_FIELD_SIZE+TYPE_FIELD_SIZE], uint16(mt))
	payloadBuf := (msgbuf)[PREPAYLOAD_SIZE:PREPAYLOAD_SIZE]
	payloadCap := cap(payloadBuf)
	payloadBuf, err := MSG_PACKER.PackMsg(msg, payloadBuf)
	if err != nil {
		return err
	}

	payloadLen := len(payloadBuf)
	if payloadLen > payloadCap {
		// exceed payload
		return fmt.Errorf("MessageConnection: message paylaod too large(%d): %v", payloadLen, msg)
	}

	var pktSize uint32 = uint32(payloadLen + PREPAYLOAD_SIZE)
	NETWORK_ENDIAN.PutUint32((msgbuf)[:SIZE_FIELD_SIZE], pktSize)
	err = mc.SendAll((msgbuf)[:pktSize])
	log.Debugf("Send message: size=%v, type=%v: %v, error=%v", pktSize, mt, msg, err)
	return err
}

// Send msg to another String through dispatcher
// Message format: [size*4B][stringID][type*2B][payload*NB]
func (mc MessageConnection) SendRelayMsg(targetStringID string, mt MsgType_t, msg interface{}) error {
	var msgbuf [MAX_MESSAGE_SIZE]byte
	copy(msgbuf[SIZE_FIELD_SIZE:SIZE_FIELD_SIZE+STRING_ID_SIZE], []byte(targetStringID))
	NETWORK_ENDIAN.PutUint16((msgbuf)[SIZE_FIELD_SIZE:SIZE_FIELD_SIZE+TYPE_FIELD_SIZE], uint16(mt))
	payloadBuf := (msgbuf)[PREPAYLOAD_SIZE:PREPAYLOAD_SIZE]
	payloadCap := cap(payloadBuf)
	payloadBuf, err := MSG_PACKER.PackMsg(msg, payloadBuf)
	if err != nil {
		return err
	}

	payloadLen := len(payloadBuf)
	if payloadLen > payloadCap {
		// exceed payload
		return fmt.Errorf("MessageConnection: message paylaod too large(%d): %v", payloadLen, msg)
	}

	var pktSize uint32 = uint32(payloadLen + PREPAYLOAD_SIZE)
	NETWORK_ENDIAN.PutUint32((msgbuf)[:SIZE_FIELD_SIZE], pktSize)
	err = mc.SendAll((msgbuf)[:pktSize])
	log.Debugf("Send message: size=%v, type=%v: %v, error=%v", pktSize, mt, msg, err)
	return err
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
	err = mc.RecvAll((*msgbuf)[SIZE_FIELD_SIZE:pktSize])
	if err != nil {
		msgbufpool.PutMsgBuf(msgbuf) // put it back on error
		return nil
	}

	var msgtype MsgType_t
	msgtype = MsgType_t(NETWORK_ENDIAN.Uint16((*msgbuf)[SIZE_FIELD_SIZE : SIZE_FIELD_SIZE+TYPE_FIELD_SIZE]))
	//msgPacker.UnpackMsg((*msgbuf)[MESSAGE_PREPAYLOAD_SIZE:pktSize),

	pinfo.Msgbuf = msgbuf
	pinfo.MsgType = msgtype
	pinfo.Payload = (*msgbuf)[PREPAYLOAD_SIZE:pktSize]

	return nil
}
