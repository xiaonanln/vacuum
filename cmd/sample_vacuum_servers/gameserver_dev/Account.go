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
	if password == "123456" {

	} else {
	}
}
