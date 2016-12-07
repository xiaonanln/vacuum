package gameserver

import "github.com/xiaonanln/vacuum/vlog"

var (
	spaceDelegate ISpaceDelegate = &SpaceDelegate{}
)

type ISpaceDelegate interface {
	OnReady(space *GSSpace)
}

func SetSpaceDelegate(delegate ISpaceDelegate) {
	spaceDelegate = delegate
}

type SpaceDelegate struct {
}

func (delegate *SpaceDelegate) OnReady(space *GSSpace) {
	vlog.Debug("%s.OnReady ...", space)
}
