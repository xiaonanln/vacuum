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
	vlog.Debug("%s GOT CLIENT !!!!!!!!!!!!!!!!", a)

	space := spaceManager.GetSpace(1)
	a.Entity.EnterSpace(space)
}
