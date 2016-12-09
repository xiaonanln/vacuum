package main

import (
	"time"

	"math/rand"

	"github.com/Sirupsen/logrus"
	"github.com/xiaonanln/typeconv"
	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/vacuum_server"
	"github.com/xiaonanln/vacuum/vlog"
)

const (
	N                 = 100
	NSERVERS          = 1
	SEND_MSG_INTERVAL = 10000 * time.Microsecond
)

type MigrateTester struct {
	vacuum.String
	val int64
}

func (pt *MigrateTester) Init() {
	vlog.Debug("!!! MigrateTester.Init %v", pt.Args())
	pt.val = 0
	pt.DeclareService("MigrateTester")
}

func (pt *MigrateTester) Loop(msg common.StringMessage) {
	pt.val += 1
	vlog.Debug("!!! MigrateTester.Loop msg %v val %v", msg, pt.val)
	msgint := typeconv.Int(msg)
	if pt.val != msgint {
		vlog.Panicf("Val is %v, but msg is %v", pt.val, msg)
	}
	pt.Migrate(1 + rand.Intn(NSERVERS))
}

func (pt *MigrateTester) Fini() {
	vlog.Info("##################################################################################################################################")
	vlog.Info("#################################################### MigrateTester.Fini %v ####################################################", pt.val)
	vlog.Info("##################################################################################################################################")
}

func (pt *MigrateTester) GetPersistentData() map[string]interface{} {
	return map[string]interface{}{
		"val": pt.val,
	}
}

func (pt *MigrateTester) LoadPersistentData(data map[string]interface{}) {
	pt.val = typeconv.Int(data["val"])
	logrus.Printf("LoadPersistentData %v", pt.val)
}

func Main() {
	for serverID := 1; serverID <= NSERVERS; serverID += 1 {
		vacuum.WaitServerReady(serverID)
	}

	stringID := vacuum.CreateString("MigrateTester", vacuum_server.ServerID())
	//vacuum.WaitServiceReady("MigrateTester", 1)
	time.Sleep(time.Second)

	for i := 0; i < N; i++ {
		vacuum.Send(stringID, i+1)
		time.Sleep(SEND_MSG_INTERVAL)
	}

	vacuum.Send(stringID, nil)

	vacuum.WaitServiceGone("MigrateTester")
	// wait for the strings to complete
}

func main() {
	vacuum.RegisterMain(Main)
	vacuum.RegisterString("MigrateTester", &MigrateTester{})
	vacuum_server.RunServer()
}
