package vacuum_server

import (
	"time"

	"github.com/xiaonanln/vacuum/config"
	"github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"
)

const (
	DISPATCHER_ADDR = ":"
)

func init() {
	// initializing the vacuum server
	config.LoadConfig()
	dispatcher_client.Initialize(1)
}

func RunServer() {
	for {
		time.Sleep(time.Second)
	}
}
