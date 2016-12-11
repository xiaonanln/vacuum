package main

import (
	"time"

	"math/rand"

	"github.com/xiaonanln/vacuum/ext/gameserver"
	_ "github.com/xiaonanln/vacuum/ext/gameserver"
	"github.com/xiaonanln/vacuum/vlog"
)

const (
	MONSTER = 1 + iota
)
const (
	NMONSTERS = 100
)

type MySpaceDelegate struct {
	gameserver.SpaceDelegate
}

func (delegate *MySpaceDelegate) OnReady(space *gameserver.GSSpace) {
	for i := 0; i < NMONSTERS; i++ {
		space.CreateEntity(MONSTER, gameserver.Vec3{100, 100, 100})
	}
}

type MyEntityDelegate struct {
	gameserver.EntityDelegate
}

func (delegate *MyEntityDelegate) OnEnterSpace(entity *gameserver.GSEntity, space *gameserver.GSSpace) {
	vlog.Debug("%s.OnEnterSpace %s, entity count %d", entity, space, space.GetEntityCount())
	entity.SetAOIDistance(100)
	if space.GetEntityCount() >= NMONSTERS {
		delegate.onAllMonstersCreated(space)
	}
}

func (delegate *MyEntityDelegate) onAllMonstersCreated(space *gameserver.GSSpace) {
	space.AddTimer(time.Second, func() {
		vlog.Info("space %s ticking", space)

		for entity, _ := range space.Entities() {
			x := gameserver.Len_t(rand.Intn(1000))
			y := gameserver.Len_t(rand.Intn(1000))
			z := gameserver.Len_t(rand.Intn(1000))
			entity.SetPos(gameserver.Vec3{x, y, z})
		}
	})
}

func main() {
	vlog.Info("gameserver_dev starting ...")
	//gameserver.GSEntity{}
	gameserver.SetSpaceDelegate(&MySpaceDelegate{})
	gameserver.SetEntityDelegate(&MyEntityDelegate{})
	gameserver.RunServer()
}
