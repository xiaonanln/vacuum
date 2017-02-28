package main

import (
	"time"

	"math/rand"

	. "github.com/xiaonanln/vacuum/ext/gameserver"
	"github.com/xiaonanln/vacuum/vlog"
)

const (
	MONSTER = 1 + iota
)
const (
	NMONSTERS = 100
)

type MySpaceDelegate struct {
	SpaceDelegate
}

func (delegate *MySpaceDelegate) OnReady(space *GSSpace) {
	vlog.Debug("%s.OnReady: kind=%v", space, space.Kind)
	if space.Kind == 0 {
		delegate.onNullSpaceReady(space)
		return
	}

	//// normal space
	//for i := 0; i < NMONSTERS; i++ {
	//	space.CreateEntity(MONSTER, Vec3{100, 100, 100})
	//}
}

func (delegate *MySpaceDelegate) onNullSpaceReady(space *GSSpace) {

}

type MyEntityDelegate struct {
	EntityDelegate
}

func (delegate *MyEntityDelegate) OnEnterSpace(entity *GSEntity, space *GSSpace) {
	vlog.Debug("%s.OnEnterSpace %s, entity count %d", entity, space, space.GetEntityCount())
	entity.SetAOIDistance(100)
	if space.GetEntityCount() >= NMONSTERS {
		delegate.onAllMonstersCreated(space)
	}
}

func (delegate *MyEntityDelegate) onAllMonstersCreated(space *GSSpace) {
	space.AddTimer(time.Second, func() {
		vlog.Info("space %s ticking", space)

		for entity, _ := range space.Entities() {
			x := Len_t(rand.Intn(1000))
			y := Len_t(rand.Intn(1000))
			z := Len_t(rand.Intn(1000))
			entity.SetPos(Vec3{x, y, z})
		}
	})
}

type Monster struct {
	GSEntityKind
}

func main() {
	vlog.Info("gameserver_dev starting ...")
	//GSEntity{}
	RegisterGSEntityKind("Account", &Account{})
	RegisterGSEntityKind("Avatar", &Avatar{})
	RegisterGSEntityKind("Monster", &Monster{})
	SetSpaceDelegate(&MySpaceDelegate{})
	SetEntityDelegate(&MyEntityDelegate{})
	RunServer()
}
