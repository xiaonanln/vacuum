package gameserver

import (
	"fmt"

	"github.com/xiaonanln/typeconv"
	"github.com/xiaonanln/vacuum/ext/entity"
	"github.com/xiaonanln/vacuum/vlog"
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
	entity.Kind = int(entityKind)
	vlog.Debug("%s.Init: kind=%v", entity, entity.Kind)
}
