package vacuum_server

import (
	"time"

	"github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"
)

const (
	DISPATCHER_ADDR = ":"
)

var ()

func RunServer() {
	dispatcher_client.RegisterVacuumServer()
	for {
		time.Sleep(time.Second)
	}
}
