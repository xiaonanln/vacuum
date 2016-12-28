package gameserver

import (
	"net"

	"github.com/xiaonanln/vacuum/proto"
	"github.com/xiaonanln/vacuum/uuid"
	"github.com/xiaonanln/vacuum/vlog"
)

type GSClient struct {
	proto.MessageConnection
	ClientID string
}

func newGSClient(conn net.Conn) *GSClient {
	return &GSClient{
		MessageConnection: proto.NewMessageConnection(conn),
		ClientID:          uuid.GenUUID(),
	}
}

func (client *GSClient) serve() {
	defer client.onServeRoutineExit()

	var err error
	for {
		err = client.RecvMsg(client)
		if err != nil {
			break
		}
	}
}

func (client *GSClient) onServeRoutineExit() {
	err := recover()
	vlog.Info("Gate client quit: %s, error=%v", client, err)
}

func (client *GSClient) HandleMsg(msg *proto.Message, pktSize uint32, msgType proto.MsgType_t) error {
	vlog.Debug("HandleMsg: pitSize=%v, msgType=%v", pktSize, msgType)
	payload := msg[proto.PREPAYLOAD_SIZE:pktSize]

	if msgType == CLIENT_RPC { // CLIENT_RPC
		client.handleClientRPC(payload)
	}
	return nil
}

func (client *GSClient) HandleRelayMsg(msg *proto.Message, pktSize uint32, targetID string) error {
	vlog.Debug("HandleRelayMsg: pktSize=%v, targetID=%v", pktSize, targetID)
	return nil
}

// RPC call from client
func (client *GSClient) handleClientRPC(payload []byte) error {
	var msg ClientRPCMessage
	if err := CLIENT_MSG_PACKER.UnpackMsg(payload, &msg); err != nil {
		return err
	}

	vlog.Debug("RPC CALL: %v", msg)
	msg.EntityID.callGSRPC(msg.Method, msg.Arguments)
	return nil
}

func (client *GSClient) clientCreateEntity(entityKind int, entityID GSEntityID) error {
	msg := ClientCreateEntityMessage{
		EntityID:   entityID,
		EntityKind: entityKind,
	}
	return client.SendMsgEx(CLIENT_CREATE_ENTITY_MESSAGE, &msg, CLIENT_MSG_PACKER)
}
