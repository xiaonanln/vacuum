package gameserver

import (
	"fmt"

	"github.com/xiaonanln/vacuum/ext/entity"
	"github.com/xiaonanln/vacuum/vlog"
	typeconv "gopkg.in/xiaonanln/typeconv.v0"
)

type GSEntity struct {
	entity.Entity
	Kind int
	Pos  Pos
}

func (entity *GSEntity) String() string {
	return fmt.Sprintf("GSEntity|%d|%s", entity.Kind, entity.ID)
}

func (entity *GSEntity) Init() {
	entityKind := typeconv.Int(entity.Args()[0])
	spaceID := SpaceID(typeconv.String(entity.Args()[1]))
	entity.Kind = int(entityKind)

	space := spaceID.getLocalSpace()
	vlog.Debug("%s.Init: space=%s", entity, space)
	space.onEntityCreated(entity)

	entityDelegate.OnReady(entity)
}
