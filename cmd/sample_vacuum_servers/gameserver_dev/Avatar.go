package main

import (
	"time"

	. "github.com/xiaonanln/vacuum/ext/gameserver"
	"github.com/xiaonanln/vacuum/vlog"
)

type Avatar struct {
	GSEntityKind
}

// called when client logined
func (a *Avatar) OnGetClient() {
	vlog.Info("%s GOT CLIENT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", a)
	a.tryEnterSpace(1)
}
func (a *Avatar) tryEnterSpace(kind int) {
	spaceID := spaceManager.GetSpace(kind)
	if spaceID != "" {
		a.EnterSpace(spaceID)
	} else {
		vlog.Info("%s.tryEnterSpace: Space %d is not ready, waiting ...", a, kind)
		a.AddCallback(time.Millisecond*100, func() {
			a.tryEnterSpace(kind) // retry
		})
	}
}

func (a *Avatar) OnDestroy() {
	vlog.Info("%s.OnDestroy ...", a)
}

func (a *Avatar) OnLoseClient() {
	vlog.Info("%s LOSE CLIENT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", a)
	a.Destroy()
}

func (a *Avatar) OnEnterSpace() {
	vlog.Info("%s ENTER SPACE %s", a, a.Space)
}
