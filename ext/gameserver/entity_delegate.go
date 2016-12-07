package gameserver

import "github.com/xiaonanln/vacuum/vlog"

var (
	entityDelegate IEntityDelegate = &EntityDelegate{}
)

type IEntityDelegate interface {
	OnReady(entity *GSEntity)
}

func SetEntityDelegate(delegate IEntityDelegate) {
	entityDelegate = delegate
}

type EntityDelegate struct {
}

func (delegate *EntityDelegate) OnReady(entity *GSEntity) {
	vlog.Debug("%s.OnReady ...", entity)
}
