package main

import (
	"strconv"

	"github.com/xiaonanln/vacuum/cmd/sample_vacuum_servers/gameserver_dev/kvdb"
	. "github.com/xiaonanln/vacuum/ext/gameserver"
	"github.com/xiaonanln/vacuum/vlog"
)

func init() {
	bootCount, _ := strconv.Atoi(kvdb.Get("bootCount", "0"))
	bootCount += 1
	vlog.Debug("BOOT COUNT: %d", bootCount)
	kvdb.Set("bootCount", strconv.Itoa(bootCount))

}

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

	avatarID := GetNilSpace().CreateEntity("Avatar", Vec3{})

	vlog.Debug("%s.Login: create Avatar %s", a, avatarID)
	a.Entity.GiveClientTo(avatarID)

}
