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
	N                 = 10000
	NSERVERS          = 2
	SEND_MSG_INTERVAL = 100 * time.Microsecond
)

type MigrateTester struct {
	val int64
}

func (pt *MigrateTester) Init(s *vacuum.String, args ...interface{}) {
	vlog.Debug("!!! MigrateTester.Init")
	pt.val = 0
	s.DeclareService("MigrateTester")
}

func (pt *MigrateTester) Loop(s *vacuum.String, msg common.StringMessage) {
	pt.val += 1
	vlog.Debug("!!! MigrateTester.Loop msg %v val %v", msg, pt.val)
	if pt.val != typeconv.Int(msg) {
		vlog.Panicf("Val is %v, but msg is %v", pt.val, msg)
	}
	s.Migrate(1 + rand.Intn(NSERVERS))
}

func (pt *MigrateTester) Fini(s *vacuum.String) {
	vlog.Info("##################################################################################################################################")
	vlog.Info("################################################## MigrateTester.Fini %v ##################################################", pt.val)
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

func Main(s *vacuum.String) {
	stringID := vacuum.CreateString("MigrateTester")
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
	for serverID := 1; serverID <= NSERVERS; serverID += 1 {
		vacuum.WaitServerReady(serverID)
	}
	vacuum.RegisterMain(Main)
	vacuum.RegisterString("MigrateTester", func() vacuum.StringDelegate {
		pt := &MigrateTester{}
		_ = vacuum.PersistentString(pt)
		return pt
	})
	vacuum_server.RunServer()
}
