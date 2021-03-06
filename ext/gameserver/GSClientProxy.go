package gameserver

import (
	"fmt"

	"github.com/xiaonanln/typeconv"
	"github.com/xiaonanln/vacuum/ext/entity"
)

type GSClientProxy struct {
	GateID   GSGateID
	ClientID GSClientID
}

func newGSClientProxy(gateID GSGateID, clientID GSClientID) *GSClientProxy {
	return &GSClientProxy{
		GateID:   gateID,
		ClientID: clientID,
	}
}

func (cp *GSClientProxy) String() string {
	return fmt.Sprintf("GSClientProxy<%s@%s>", cp.ClientID, cp.GateID)
}

func (cp *GSClientProxy) callClient(entityID GSEntityID, methodName string, args []interface{}) {
	// call the gate
	entity.EntityID(cp.GateID).Call("CallClient", cp.ClientID, entityID, methodName, args)
}

// notify the client to change owner
func (cp *GSClientProxy) notifyChangeOwner(ownerID GSEntityID, otherID GSEntityID, otherKindName string) {
	entity.EntityID(cp.GateID).Call("NotifyClientChangeOwner", cp.ClientID, ownerID, otherID, otherKindName)
}

func (cp *GSClientProxy) getClientProxyData() interface{} {
	if cp != nil {
		return []interface{}{cp.GateID, cp.ClientID}
	} else {
		return nil
	}
}

func (cp *GSClientProxy) setClientProxyData(data interface{}) {
	datalist := data.([]interface{})
	cp.GateID = GSGateID(typeconv.String(datalist[0]))
	cp.ClientID = GSClientID(typeconv.String(datalist[1]))
}
