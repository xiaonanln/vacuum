package proto

import (
	"net"

	"encoding/binary"

	"fmt"

	log "github.com/Sirupsen/logrus"

	"sync"

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

	RELAY_MASK = 0x80000000
)

var (
	NETWORK_ENDIAN = binary.LittleEndian
	messagePool    = sync.Pool{
		New: func() interface{} {
			return &Message{}
		},
	}
)

type MessageConnection struct {
	netutil.BinaryConnection
}

func NewMessageConnection(conn net.Conn) MessageConnection {
	return MessageConnection{BinaryConnection: netutil.NewBinaryConnection(conn)}
}

type Message [MAX_MESSAGE_SIZE]byte

func allocMessage() *Message {
	return messagePool.Get().(*Message)
}

func (m *Message) release() {
	messagePool.Put(m)
}

// Send msg to/from dispatcher
// Message format: [size*4B][type*2B][payload*NB]
func (mc MessageConnection) SendMsg(mt MsgType_t, msg interface{}) error {
	msgbuf := allocMessage()

	NETWORK_ENDIAN.PutUint16((msgbuf)[SIZE_FIELD_SIZE:SIZE_FIELD_SIZE+TYPE_FIELD_SIZE], uint16(mt))
	payloadBuf := (msgbuf)[PREPAYLOAD_SIZE:PREPAYLOAD_SIZE]
	payloadCap := cap(payloadBuf)
	payloadBuf, err := MSG_PACKER.PackMsg(msg, payloadBuf)
	if err != nil {
		msgbuf.release()
		return err
	}

	payloadLen := len(payloadBuf)
	if payloadLen > payloadCap {
		// exceed payload
		msgbuf.release()
		return fmt.Errorf("MessageConnection: message paylaod too large(%d): %v", payloadLen, msg)
	}

	var pktSize uint32 = uint32(payloadLen + PREPAYLOAD_SIZE)
	NETWORK_ENDIAN.PutUint32((msgbuf)[:SIZE_FIELD_SIZE], pktSize)
	err = mc.SendAll((msgbuf)[:pktSize])
	msgbuf.release()
	log.Debugf("Send message: size=%v, type=%v: %v, error=%v", pktSize, mt, msg, err)
	return err
}

// Send msg to another String through dispatcher
// Message format: [size*4B][stringID][type*2B][payload*NB]
func (mc MessageConnection) SendRelayMsg(targetID string, mt MsgType_t, msg interface{}) error {
	msgbuf := allocMessage()
	copy(msgbuf[SIZE_FIELD_SIZE:SIZE_FIELD_SIZE+STRING_ID_SIZE], []byte(targetID))

	NETWORK_ENDIAN.PutUint16((msgbuf)[SIZE_FIELD_SIZE+STRING_ID_SIZE:SIZE_FIELD_SIZE+STRING_ID_SIZE+TYPE_FIELD_SIZE], uint16(mt))
	payloadBuf := (msgbuf)[RELAY_PREPAYLOAD_SIZE:RELAY_PREPAYLOAD_SIZE]
	payloadCap := cap(payloadBuf)
	payloadBuf, err := MSG_PACKER.PackMsg(msg, payloadBuf)
	if err != nil {
		msgbuf.release()
		return err
	}

	payloadLen := len(payloadBuf)
	if payloadLen > payloadCap {
		// exceed payload
		msgbuf.release()
		return fmt.Errorf("MessageConnection: message paylaod too large(%d): %v", payloadLen, msg)
	}

	var pktSize uint32 = uint32(payloadLen + RELAY_PREPAYLOAD_SIZE)
	NETWORK_ENDIAN.PutUint32((msgbuf)[:SIZE_FIELD_SIZE], pktSize|RELAY_MASK) // set highest bit of size to 1 to indicate a relay msg
	err = mc.SendAll((msgbuf)[:pktSize])
	msgbuf.release()
	log.Debugf("Send relay message: size=%v, targetID=%s, type=%v: %v, error=%v", pktSize, targetID, mt, msg, err)
	return err
}

type MessageHandler interface {
	HandleMsg(msg *Message, pktSize uint32, msgType MsgType_t) error
	HandleRelayMsg(msg *Message, pktSize uint32, targetID string) error
}

func (mc MessageConnection) RecvMsg(handler MessageHandler) error {
	msg := allocMessage()

	pktSizeBuf := msg[:SIZE_FIELD_SIZE]
	err := mc.RecvAll(pktSizeBuf)
	if err != nil {
		return err
	}

	var pktSize uint32 = NETWORK_ENDIAN.Uint32(pktSizeBuf)
	isRelayMsg := false

	if pktSize&RELAY_MASK != 0 {
		// this is a relay msg
		isRelayMsg = true
		pktSize -= RELAY_MASK
	}

	log.Debugf("RecvMsg: pktsize=%v, isRelayMsg=%v", pktSize, isRelayMsg)

	if pktSize > MAX_MESSAGE_SIZE {
		// pkt size is too large
		msg.release()
		return fmt.Errorf("message packet too large: %v", pktSize)
	}

	err = mc.RecvAll((msg)[SIZE_FIELD_SIZE:pktSize]) // receive the msg type and payload
	if err != nil {
		msg.release()
		return err
	}

	log.WithFields(log.Fields{"pktSize": pktSize, "isRelayMsg": isRelayMsg}).Debugf("RecvMsg")
	if isRelayMsg {
		// if it is a relay msg, we just relay what we receive without interpret the payload
		targetID := string(msg[SIZE_FIELD_SIZE : SIZE_FIELD_SIZE+STRING_ID_SIZE])
		err = handler.HandleRelayMsg(msg, pktSize, targetID)
	} else {
		var msgtype MsgType_t
		msgtype = MsgType_t(NETWORK_ENDIAN.Uint16((msg)[SIZE_FIELD_SIZE : SIZE_FIELD_SIZE+TYPE_FIELD_SIZE]))
		err = handler.HandleMsg(msg, pktSize, msgtype)
	}

	msg.release()
	return err
}
