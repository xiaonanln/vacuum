package gs_space

import (
	"github.com/xiaonanln/typeconv"
	"github.com/xiaonanln/vacuum/ext/entity"
	"github.com/xiaonanln/vacuum/vlog"
)

type GSSpace struct {
	entity.Entity
}

func (space *GSSpace) Init() {
	spaceKind := typeconv.Int(space.Args()[0])
	vlog.Info("GSSpace.Init kind=%v", spaceKind)
}
