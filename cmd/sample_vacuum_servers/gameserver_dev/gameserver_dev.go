package main

import (
	"time"

	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/ext/gameserver"
	_ "github.com/xiaonanln/vacuum/ext/gameserver"
	"github.com/xiaonanln/vacuum/vacuum_server"
	"github.com/xiaonanln/vacuum/vlog"
)

func main() {
	vlog.Info("gameserver_dev starting ...")
	//gameserver.GSEntity{}
	vacuum.RegisterMain(func() {
		spaceID := gameserver.CreateSpace(0)
		vlog.Info("Create space: %s", spaceID)
		time.Sleep(3 * time.Second)
	})
	vacuum_server.RunServer()
}
