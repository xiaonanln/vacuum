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
	N                 = 1000
	NSERVERS          = 2
	SEND_MSG_INTERVAL = 10 * time.Microsecond
)

type MigrateTester struct {
	val int64
}

func (pt *MigrateTester) Init(s *vacuum.String) {
	vlog.Debug("!!! MigrateTester.Init %v", s.Args())
	pt.val = 0
	s.DeclareService("MigrateTester")
}

func (pt *MigrateTester) Loop(s *vacuum.String, msg common.StringMessage) {
	pt.val += 1
	vlog.Debug("!!! MigrateTester.Loop msg %v val %v", msg, pt.val)
	msgint := typeconv.Int(msg)
	if pt.val != msgint {
		vlog.Panicf("Val is %v, but msg is %v", pt.val, msg)
	}
	if msgint < N {
		s.Migrate(1 + rand.Intn(NSERVERS))
	} else {
		s.Migrate(int(msgint))
	}
}

func (pt *MigrateTester) Fini(s *vacuum.String) {
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

func Main(s *vacuum.String) {
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
