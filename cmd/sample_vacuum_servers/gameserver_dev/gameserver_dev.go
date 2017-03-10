package main

import (
	. "github.com/xiaonanln/vacuum/ext/gameserver"
	"github.com/xiaonanln/vacuum/vlog"
)

const (
	MONSTER = 1 + iota
)
const (
	NMONSTERS = 100
)

type Monster struct {
	GSEntityKind
}

func main() {
	vlog.Info("gameserver_dev starting ...")
	//GSEntity{}
	RegisterGSEntityKind("Account", &Account{})
	RegisterGSEntityKind("Avatar", &Avatar{})
	RegisterGSEntityKind("Monster", &Monster{})
	SetSpaceDelegate(&MySpaceDelegate{})
	RunServer()
}
