package main

import (
	"time"

	"gopkg.in/xiaonanln/typeconv.v0"

	"os"

	"github.com/Sirupsen/logrus"
	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/vacuum_server"
)

type PersistentTester struct {
	vacuum.String
	val int64
}

func (pt *PersistentTester) Init() {
	pt.val = 0
	pt.DeclareService("PersistentTester")
}

func (pt *PersistentTester) Loop(msg common.StringMessage) {
	pt.val += typeconv.Int(msg)
	pt.Save()
}

func (pt *PersistentTester) Fini() {
	logrus.Printf("!!!!!!!!!!! Fini %v!!!!!!!!!!!!!!", pt.val)
}

func (pt *PersistentTester) IsPersistent() bool {
	return true
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

type Main struct {
	vacuum.String
}

func (s *Main) Init() {
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
	os.Exit(0)
}

func main() {
	vacuum.RegisterString("Main", &Main{})
	vacuum.RegisterString("PersistentTester", &PersistentTester{})
	vacuum_server.RunServer()
}
