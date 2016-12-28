package gameserver

import "github.com/xiaonanln/vacuum/proto"

const (
	CLIENT_RPC                   = 1 + iota
	CLIENT_CREATE_ENTITY_MESSAGE = 1 + iota
)

type ClientRPCMessage struct {
	EntityID  GSEntityID    `json:"E"`
	Method    string        `json:"M"`
	Arguments []interface{} `json:"A"`
}

type ClientCreateEntityMessage struct {
	EntityKind int        `json:"K"`
	EntityID   GSEntityID `json:"E"`
}

var (
	CLIENT_MSG_PACKER = proto.JSONMsgPacker{}
)