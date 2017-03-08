package gameserver

import (
	"fmt"

	"net"

	"log"

	"sync"

	"github.com/xiaonanln/typeconv"
	"github.com/xiaonanln/vacuum/ext/entity"
	"github.com/xiaonanln/vacuum/netutil"
	"github.com/xiaonanln/vacuum/vlog"
)

type GSGateID entity.EntityID

type GSGate struct {
	entity.Entity
	ID          GSGateID
	clientsLock sync.RWMutex
	clients     map[GSClientID]*GSClient
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
	// called from GSGate.Entity routine
	gate.ID = GSGateID(gate.Entity.ID)

	// initialize clients
	gate.clients = map[GSClientID]*GSClient{}
}

func (gate *GSGate) OnReady() {
	// start goroutine to serve clients
	gateIndex := typeconv.Int(gate.Args()[0])
	port := typeconv.Int(gate.Args()[1])
	vlog.Debug("Initializing gate %v port %v: %s ...", gateIndex, port, gate)

	serveAddr := fmt.Sprintf(":%d", port)
	go netutil.ServeTCPForever(serveAddr, gate)
}

// handle client connection to gate
func (gate *GSGate) ServeTCPConnection(conn net.Conn) {
	// called from TCP Service
	vlog.Debug("%s: new connection %s ...", gate, conn.RemoteAddr())
	client := newGSClient(gate, conn)
	gate.clientsLock.Lock()
	gate.clients[client.ClientID] = client // add client to gate clients-map
	gate.clientsLock.Unlock()

	bootEntityKindName := gameserverConfig.BootEntityKind
	entityID := nilSpace.CreateEntity(bootEntityKindName, Vec3{})
	// set entity client

	client.setOwner(entityID)
	client.letClientCreateEntity(bootEntityKindName, entityID) // create entity on client side
	entityID.notifyGetClient(gate.ID, client.ClientID)

	go client.serve()
}

func (gate *GSGate) CallClient(clientID GSClientID, entityID GSEntityID, methodName string, args []interface{}) {
	// Send RPC call to the client, called from GSGate.Entity routine
	gate.clientsLock.RLock()
	client, ok := gate.clients[clientID]
	gate.clientsLock.RUnlock()
	if !ok {
		log.Panicf("%s.CallClient: %s: Client %s not found", gate, methodName, clientID)
	}

	client.clientCallEntityMethod(entityID, methodName, args)
}

// notify the client to change owner
func (gate *GSGate) NotifyClientChangeOwner(clientID GSClientID, ownerID GSEntityID, otherID GSEntityID, otherKindName string) {
	gate.clientsLock.RLock()
	client, ok := gate.clients[clientID]
	gate.clientsLock.RUnlock()
	if !ok {
		log.Panicf("%s.NotifyClientChangeOwner: client %s not found", gate, clientID)
		return
	}

	// todo: optimize client operations
	client.notifyChangeOwner(ownerID, otherID)
	client.letClientDestroyEntity(ownerID)
	client.letClientCreateEntity(otherKindName, otherID)
	// tell other entity to get the client
	otherID.notifyGetClient(gate.ID, client.ClientID)
}

func (gate *GSGate) onClientDisconnect(client *GSClient) {
	// called from GSClient serve routine
	gate.clientsLock.Lock()
	delete(gate.clients, client.ClientID)
	gate.clientsLock.Unlock()
}
