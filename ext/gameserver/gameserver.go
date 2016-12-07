package gameserver

import (
	"github.com/xiaonanln/vacuum/ext/entity"
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
}
