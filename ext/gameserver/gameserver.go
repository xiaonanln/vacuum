package gameserver

import (
	"time"

	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/ext/entity"
	"github.com/xiaonanln/vacuum/vacuum_server"
	"github.com/xiaonanln/vacuum/vlog"
)

const (
	DEFAULT_GATES_NUM = 1
)

func init() {
	vlog.Debug("Register gameserver entities ...")
	entity.RegisterEntity("GSSpace", &GSSpace{})
	entity.RegisterEntity("GSEntity", &GSEntity{})
	entity.RegisterEntity("GSGate", &GSGate{})

}

func gameserverMain() {
	gameserverConfig := loadGameserverConfig()
	vlog.Debug("Gameserver config: %v", gameserverConfig)

	runGates(gameserverConfig)

	time.Sleep(10 * time.Second)
}

func RunServer() {
	vacuum.RegisterMain(gameserverMain)
	vacuum_server.RunServer()
}
