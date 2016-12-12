package gameserver

import (
	"time"

	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/config"
	"github.com/xiaonanln/vacuum/ext/entity"
	"github.com/xiaonanln/vacuum/ext/gameserver"
	"github.com/xiaonanln/vacuum/vacuum_server"
	"github.com/xiaonanln/vacuum/vlog"
)

const (
	SPACE_ENTITY_TYPE = "GSSpace"
	ENTITY_TYPE       = "GSEntity"
)

func init() {
	vlog.Debug("Register gameserver entities ...")
	entity.RegisterEntity(SPACE_ENTITY_TYPE, &GSSpace{})
	entity.RegisterEntity(ENTITY_TYPE, &GSEntity{})
	entity.RegisterEntity("GSGate", &GSGate{})

}

func RunServer() {
	vacuum.RegisterMain(func() {
		gameserverConfig := config.LoadExtraConfig("gameserver")
		vlog.Debug("Gameserver config: %v", gameserverConfig)

		spaceID := gameserver.CreateSpace(0)
		vlog.Info("Create space: %s", spaceID)
		time.Sleep(10 * time.Second)
	})
	vacuum_server.RunServer()
}
