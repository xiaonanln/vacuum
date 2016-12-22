package gameserver

import "github.com/xiaonanln/vacuum/proto"

const (
	CLIENT_RPC = 1 + iota
)

type ClientRPCMessage struct {
	Method    string        `json:"M"`
	Arguments []interface{} `json:"A"`
}

var (
	CLIENT_MSG_PACKER = proto.JSONMsgPacker{}
)
