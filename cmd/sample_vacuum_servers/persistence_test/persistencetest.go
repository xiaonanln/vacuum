package main

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/xiaonanln/typeconv"
	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/vacuum_server"
)

type PersistentTester struct {
	val int64
}

func (pt *PersistentTester) Init(s *vacuum.String) {
	pt.val = 0
	s.DeclareService("PersistentTester")
}

func (pt *PersistentTester) Loop(s *vacuum.String, msg common.StringMessage) {
	pt.val += typeconv.Int(msg)
	s.Save()
}

func (pt *PersistentTester) Fini(s *vacuum.String) {
	logrus.Printf("!!!!!!!!!!! Fini %v!!!!!!!!!!!!!!", pt.val)
}

func (pt *PersistentTester) GetPersistentData() map[string]interface{} {
	return map[string]interface{}{
		"val": pt.val,
	}
}

func (pt *PersistentTester) LoadPersistentData(data map[string]interface{}) {
	pt.val = typeconv.Int(data["val"])
	logrus.Printf("!!!!!!!!!!! LoadPersistentData %v!!!!!!!!!!!!!!", pt.val)
}

func Main(s *vacuum.String) {
	stringID := vacuum.CreateString("PersistentTester")
	vacuum.WaitServiceReady("PersistentTester", 1)
	vacuum.Send(stringID, 1)
	vacuum.Send(stringID, nil)
	vacuum.WaitServiceGone("PersistentTester")
	time.Sleep(time.Second)

	for {
		vacuum.LoadString("PersistentTester", stringID)
		vacuum.WaitServiceReady("PersistentTester", 1)
		vacuum.Send(stringID, 1)
		vacuum.Send(stringID, nil)
		vacuum.WaitServiceGone("PersistentTester")
		time.Sleep(time.Second)
	}

	time.Sleep(3 * time.Second)
}

func main() {
	vacuum.RegisterMain(Main)
	vacuum.RegisterString("PersistentTester", func() vacuum.StringDelegate {
		pt := &PersistentTester{}
		vacuum.AssurePersistent(pt)
		return pt
	})
	vacuum_server.RunServer()
}
