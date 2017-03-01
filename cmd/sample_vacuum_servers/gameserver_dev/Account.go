package main

import (
	. "github.com/xiaonanln/vacuum/ext/gameserver"
	"github.com/xiaonanln/vacuum/vlog"
)

type Account struct {
	GSEntityKind
}

func (a *Account) Login_OwnClient(username string, password string) {
	vlog.Info("%s.Login %s %s", a, username, password)
	if password != "123456" {
		a.Entity.CallClient("OnLogin", false)
		return
	}

	a.Entity.CallClient("OnLogin", true) // tell client that login ok
	// create the new Avatar entity

	avatarID := CreateGSEntityAnywhere("Avatar")
	vlog.Debug("%s.Login: create Avatar %s", a, avatarID)
	a.Entity.GiveClientTo(avatarID)
}
