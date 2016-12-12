package gameserver

import (
	"fmt"

	"net"

	"github.com/xiaonanln/typeconv"
	"github.com/xiaonanln/vacuum/ext/entity"
	"github.com/xiaonanln/vacuum/netutil"
	"github.com/xiaonanln/vacuum/vlog"
)

const ()

type GSGate struct {
	entity.Entity
}

func init() {
}

func runGates(config *gameserverConfig) {
	for i := 0; i < config.GatesNum; i++ {
		port := config.GatesStartPort + i*config.GatesPortStep
		entity.CreateEntity("GSGate", i, port)
	}
}

func (gate *GSGate) Init() {
	gateIndex := typeconv.Int(gate.Args()[0])
	port := typeconv.Int(gate.Args()[1])
	vlog.Debug("Initializing gate %v port %v: %s ...", gateIndex, port, gate)
	serveAddr := fmt.Sprintf(":%d", port)
	netutil.ServeTCPForever(serveAddr, &gateServerDelegate{gate: gate})
}

type gateServerDelegate struct {
	gate *GSGate
}

// handle client connection to gate
func (delegate *gateServerDelegate) ServeTCPConnection(conn net.Conn) {
	vlog.Debug("%s: new connection %s ...", delegate.gate, conn.RemoteAddr())
	client := newGSClient(conn)
	go client.serve()
}
