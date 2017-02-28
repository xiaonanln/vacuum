package gameserver

import (
	"fmt"

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
