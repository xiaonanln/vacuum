package gameserver

import (
	"github.com/xiaonanln/vacuum/ext/entity"
	"github.com/xiaonanln/vacuum/proto"
)

const (
	CLIENT_RPC = 1 + iota
)

type ClientRPCMessage struct {
	EntityID  entity.EntityID `json:"E"`
	Method    string          `json:"M"`
	Arguments []interface{}   `json:"A"`
}

var (
	CLIENT_MSG_PACKER = proto.JSONMsgPacker{}
)
