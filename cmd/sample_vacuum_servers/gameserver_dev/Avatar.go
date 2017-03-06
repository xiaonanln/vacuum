package main

import (
	"time"

	"github.com/xiaonanln/goTimer"
	. "github.com/xiaonanln/vacuum/ext/gameserver"
	"github.com/xiaonanln/vacuum/vlog"
)

type Avatar struct {
	GSEntityKind
}

// called when client logined
func (a *Avatar) OnGetClient() {
	vlog.Debug("%s GOT CLIENT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", a)

	space := spaceManager.GetSpace(1)
	a.Entity.EnterSpace(space)
	timer.AddCallback(time.Second*5, func() {
		a.Entity.Destroy()
	})
	//a.Entity.Migrate(vacuum_server.ServerID())
}
