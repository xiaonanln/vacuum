package gameserver

import (
	"fmt"

	"net"

	"log"

	"github.com/xiaonanln/typeconv"
	"github.com/xiaonanln/vacuum/ext/entity"
	"github.com/xiaonanln/vacuum/netutil"
	"github.com/xiaonanln/vacuum/vlog"
)

type GSGateID entity.EntityID

type GSGate struct {
	entity.Entity
	ID      GSGateID
	clients map[GSClientID]*GSClient
}

func init() {
}

func runGates(config *GameserverConfig) {
	for i := 0; i < config.GatesNum; i++ {
		port := config.GatesStartPort + i*config.GatesPortStep
		entity.CreateEntity("GSGate", i, port)
	}
}

func (gate *GSGate) String() string {
	return fmt.Sprintf("GSGate<%s>", gate.ID)
}

func (gate *GSGate) Init() {
	gate.ID = GSGateID(gate.Entity.ID)

	// initialize clients
	gate.clients = map[GSClientID]*GSClient{}

	gateIndex := typeconv.Int(gate.Args()[0])
	port := typeconv.Int(gate.Args()[1])
	vlog.Debug("Initializing gate %v port %v: %s ...", gateIndex, port, gate)

	// start goroutine to serve clients
	serveAddr := fmt.Sprintf(":%d", port)
	go netutil.ServeTCPForever(serveAddr, gate)
}

// handle client connection to gate
func (gate *GSGate) ServeTCPConnection(conn net.Conn) {
	vlog.Debug("%s: new connection %s ...", gate, conn.RemoteAddr())
	client := newGSClient(gate.ID, conn)
	gate.clients[client.ClientID] = client // add client to gate clients-map

	bootEntityKindName := gameserverConfig.BootEntityKind
	entityID := nilSpace.CreateEntity(bootEntityKindName, Vec3{})
	// set entity client

	client.setOwner(entityID)
	client.clientCreateEntity(bootEntityKindName, entityID) // create entity on client side
	entityID.notifyGetClient(gate.ID, client.ClientID)

	go client.serve()
}

func (gate *GSGate) CallClient(clientID GSClientID, entityID GSEntityID, methodName string, args []interface{}) {
	// Send RPC call to the client
	client, ok := gate.clients[clientID]
	if !ok {
		log.Panicf("%s.CallClient: %s: Client %s not found", gate, methodName, clientID)
	}

	client.clientCallEntityMethod(entityID, methodName, args)
}
