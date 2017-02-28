package gameserver

import "github.com/xiaonanln/vacuum/proto"

const (
	CLIENT_TO_SERVER_OWN_CLIENT_RPC = 1 + iota
	CLIENT_CREATE_ENTITY_MESSAGE    = 1 + iota
	SERVER_TO_CLIENT_RPC            = 1 + iota
)

type ClientRPCMessage struct {
	EntityID  GSEntityID    `json:"E"`
	Method    string        `json:"M"`
	Arguments []interface{} `json:"A"`
}

type ClientCreateEntityMessage struct {
	EntityKind string     `json:"K"`
	EntityID   GSEntityID `json:"E"`
}

type ServerToClientRPCMessage struct {
	EntityID  GSEntityID    `json:"E"`
	Method    string        `json:"M"`
	Arguments []interface{} `json:"A"`
}

var (
	CLIENT_MSG_PACKER = proto.JSONMsgPacker{}
)
