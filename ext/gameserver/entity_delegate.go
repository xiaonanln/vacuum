package gameserver

import "github.com/xiaonanln/vacuum/vlog"

var (
	entityDelegate IEntityDelegate = &EntityDelegate{}
)

type IEntityDelegate interface {
	OnReady(entity *GSEntity)
	OnEnterSpace(entity *GSEntity, space *GSSpace)
	OnEnterAOI(entity *GSEntity, other *GSEntity)
	OnLeaveAOI(entity *GSEntity, other *GSEntity)
}

func SetEntityDelegate(delegate IEntityDelegate) {
	entityDelegate = delegate
}

type EntityDelegate struct {
}

func (delegate *EntityDelegate) OnReady(entity *GSEntity) {
	vlog.Debug("%s.OnReady ...", entity)
}

func (delegate *EntityDelegate) OnEnterSpace(entity *GSEntity, space *GSSpace) {
	vlog.Debug("%s.OnEnterSpace %s", entity, space)
}

func (delegate *EntityDelegate) OnEnterAOI(entity *GSEntity, other *GSEntity) {
	vlog.Debug("%s.OnEnterAOI: %s", entity, other)
}

func (delegate *EntityDelegate) OnLeaveAOI(entity *GSEntity, other *GSEntity) {
	vlog.Debug("%s.OnLeaveAOI: %s", entity, other)
}
