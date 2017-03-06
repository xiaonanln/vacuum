package main

import (
	"strconv"

	"time"

	"github.com/xiaonanln/goTimer"
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

	loginingAvatarID GSEntityID
}

func (a *Account) Login_OwnClient(username string, password string) {

	vlog.Info("%s.Login %s %s", a, username, password)
	if password != "123456" {
		a.Entity.CallClient("OnLogin", false)
		return
	}

	a.Entity.CallClient("OnLogin", true) // tell client that login ok

	// get avatar id from kvdb
	var avatarID GSEntityID
	avatarID = GSEntityID(kvdb.Get("AvatarID-"+username, ""))
	if avatarID == "" {
		// new account
		avatarID = CreateGSEntityLocally("Avatar")
		vlog.Debug("%s.Login: create Avatar %s", a, avatarID)
		kvdb.Set("AvatarID-"+username, string(avatarID))
		a.loginingAvatarID = avatarID

		a.onAvatarReadyLocally(avatarID)
		return
	}

	vlog.Debug("%s.Login: loading avatar %s ...", a, avatarID)
	LoadGSEntity("Avatar", avatarID)
	a.loginingAvatarID = avatarID

	timer.AddCallback(time.Second, func() {
		a.onLoadAvatarComplete()
	})

}

func (a *Account) onLoadAvatarComplete() {
	vlog.Debug("%s.onLoadAvatarComplete ..,", a)

	a.Entity.MigrateTowards(a.loginingAvatarID)
}

func (a *Account) onAvatarReadyLocally(avatarID GSEntityID) {
	a.Entity.GiveClientTo(avatarID)
	a.Entity.Destroy()
}

func (a *Account) OnEnterSpace() {
	avatarID := a.loginingAvatarID
	vlog.Debug("%s.OnEnterSpace: login avatar = %s", a, avatarID)
	if avatarID == "" {
		return
	}

	if avatarID.GetLocalGSEntity() == nil {
		// avatar not found, ...
		vlog.Warn("%s.OnEnterSpace: avatar %s not found on local server, login failed", a, avatarID)
		a.Entity.Destroy()
		return
	}

	a.onAvatarReadyLocally(avatarID)
}
