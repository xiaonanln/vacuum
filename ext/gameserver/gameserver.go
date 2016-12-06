package gameserver

import (
	"github.com/xiaonanln/vacuum/ext/entity"
	"github.com/xiaonanln/vacuum/ext/gameserver/space"
	"github.com/xiaonanln/vacuum/vlog"
)

const (
	SPACE_ENTITY_TYPE = "GSSpace"
	ENTITY_NAME       = "GSEntity"
)

func init() {
	vlog.Debug("Register gameserver entities ...")
	entity.RegisterEntity(SPACE_ENTITY_TYPE, &gs_space.GSSpace{})
	entity.RegisterEntity(ENTITY_NAME, &GSEntity{})
}
