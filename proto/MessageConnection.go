package proto

import (
	"unsafe"

	"github.com/xiaonanln/vacuum/netutil"
)

const (
	MAX_MSG_SIZE = 64*1024 - unsafe.Sizeof(msgtype_t(0))
)

var (
	msgPacker = MessagePackMsgPacker{}
)

type MessageConnection struct {
	netutil.BinaryConnection
}

func (mc MessageConnection) SendMsg(mt msgtype_t, msg interface{}) {
	msgPacker.PackMsg(msg, buf)
}
