package main

import (
	. "github.com/xiaonanln/vacuum/ext/gameserver"
	"github.com/xiaonanln/vacuum/vlog"
)

type Avatar struct {
	GSEntityKind
}

// called when client logined
func (a *Avatar) OnGetClient() {
	vlog.Info("%s GOT CLIENT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", a)

	space := spaceManager.GetSpace(1)
	a.EnterSpace(space)
}

func (a *Avatar) OnDestroy() {
	vlog.Info("%s.OnDestroy ...", a)
}

func (a *Avatar) OnLoseClient() {
	vlog.Info("%s LOSE CLIENT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", a)
}

func (a *Avatar) OnEnterSpace() {
	vlog.Info("%s ENTER SPACE %s", a, a.Space)
}
