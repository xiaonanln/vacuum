package gameserver

import "github.com/xiaonanln/vacuum/vlog"

var (
	spaceDelegate ISpaceDelegate
)

type ISpaceDelegate interface {
	OnLoaded(space *GSSpace)
}

func SetSpaceDelegate(delegate ISpaceDelegate) {
	spaceDelegate = delegate
}

type SpaceDelegate struct {
}

func (delegate *SpaceDelegate) OnLoaded(space *GSSpace) {
	vlog.Debug("%s.OnInit ...", space)
}
