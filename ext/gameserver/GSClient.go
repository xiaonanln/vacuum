package gameserver

import (
	"net"

	"fmt"

	"sync"

	"github.com/xiaonanln/vacuum/proto"
	"github.com/xiaonanln/vacuum/uuid"
	"github.com/xiaonanln/vacuum/vlog"
)

type GSClientID string

type GSClient struct {
	proto.MessageConnection
	gate     *GSGate
	ClientID GSClientID

	sync.RWMutex
	ownerID GSEntityID
}

func newGSClient(gate *GSGate, conn net.Conn) *GSClient {
	return &GSClient{
		gate:              gate,
		MessageConnection: proto.NewMessageConnection(conn),
		ClientID:          GSClientID(uuid.GenUUID()),
	}
}

func (client *GSClient) setOwner(entityID GSEntityID) {
	client.Lock()
	client.ownerID = entityID
	client.Unlock()
}

func (client *GSClient) serve() {
	defer client.onServeRoutineExit()

	var err error
	for {
		err = client.RecvMsg(client)
		if err != nil {
			panic(err)
		}
	}
}

func (client *GSClient) onServeRoutineExit() {
	err := recover()
	vlog.Info("Gate client quit: %s, error=%v", client, err)
	// remove self from GSGate
	client.gate.onClientDisconnect(client)

	// notify the owner entity
	client.Lock()
	if client.ownerID != "" {
		ownerID := client.ownerID
		client.ownerID = ""
		client.Unlock()
		ownerID.notifyLoseClient(client.gate.ID, client.ClientID)
	} else {
		client.Unlock()
	}
}

func (client *GSClient) HandleMsg(msg *proto.Message, pktSize uint32, msgType proto.MsgType_t) error {
	vlog.Debug("HandleMsg: pktSize=%v, msgType=%v", pktSize, msgType)
	payload := msg[proto.PREPAYLOAD_SIZE:pktSize]

	switch msgType {
	case CLIENT_TO_SERVER_OWN_CLIENT_RPC:
		return client.handleClientOwnClientRPC(payload)
	default:
		return fmt.Errorf("%s: invalid message type %v", client, msgType)
	}
	return nil
}

func (client *GSClient) HandleRelayMsg(msg *proto.Message, pktSize uint32, targetID string) error {
	vlog.Debug("HandleRelayMsg: pktSize=%v, targetID=%v", pktSize, targetID)
	return nil
}

// RPC call from client
func (client *GSClient) handleClientOwnClientRPC(payload []byte) error {
	var msg ClientRPCMessage
	if err := CLIENT_MSG_PACKER.UnpackMsg(payload, &msg); err != nil {
		return err
	}

	vlog.Debug("RPC CALL: %v, OWN CLIENT", msg)
	msg.EntityID.callGSRPC_OwnClient(msg.Method, msg.Arguments)
	return nil
}

func (client *GSClient) clientCreateEntity(kindName string, entityID GSEntityID) error {
	msg := ClientCreateEntityMessage{
		EntityID:   entityID,
		EntityKind: kindName,
	}
	return client.SendMsgEx(CLIENT_CREATE_ENTITY_MESSAGE, &msg, CLIENT_MSG_PACKER)
}

func (client *GSClient) clientCallEntityMethod(entityID GSEntityID, methodName string, args []interface{}) error {
	msg := ServerToClientRPCMessage{
		EntityID:  entityID,
		Method:    methodName,
		Arguments: args,
	}
	return client.SendMsgEx(SERVER_TO_CLIENT_RPC, &msg, CLIENT_MSG_PACKER)
}

func (client *GSClient) notifyChangeOwner(ownerID GSEntityID, otherID GSEntityID) {
	client.Lock()
	if client.ownerID != ownerID {
		client.Unlock()
		vlog.Warn("%s.notifyChangeOwner: ownerID is %s, but notified by %s", client, client.ownerID, ownerID)
		return
	}

	client.ownerID = otherID
	client.Unlock()

	// tell the owner
}
