package vacuum_server

import (
	"time"

	"github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"
)

const (
	DISPATCHER_ADDR = ":"
)

func init() {
	// initializing the vacuum server
	dispatcher_client.RegisterVacuumServer(1)
}

func RunServer(serverID int) {
	for {
		time.Sleep(time.Second)
	}
}
