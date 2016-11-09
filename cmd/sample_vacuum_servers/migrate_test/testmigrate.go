package main

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/xiaonanln/typeconv"
	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/vacuum_server"
)

type MigrateTester struct {
	val int64
}

func (pt *MigrateTester) Init(s *vacuum.String, args ...interface{}) {
	pt.val = 0
	s.DeclareService("MigrateTester")
}

func (pt *MigrateTester) Loop(s *vacuum.String, msg common.StringMessage) {
	pt.val += typeconv.Int(msg)
	s.Save()
	s.Migrate(1)
}

func (pt *MigrateTester) Fini(s *vacuum.String) {
	logrus.Printf("!!!!!!!!!!! Fini %v!!!!!!!!!!!!!!", pt.val)
}

func (pt *MigrateTester) GetPersistentData() map[string]interface{} {
	return map[string]interface{}{
		"val": pt.val,
	}
}

func (pt *MigrateTester) LoadPersistentData(data map[string]interface{}) {
	pt.val = typeconv.Int(data["val"])
	logrus.Printf("!!!!!!!!!!! LoadPersistentData %v!!!!!!!!!!!!!!", pt.val)
}

func Main(s *vacuum.String) {
	stringID := vacuum.CreateString("MigrateTester")
	vacuum.WaitServiceReady("MigrateTester", 1)
	time.Sleep(time.Second)

	for {
		vacuum.Send(stringID, 1)
		time.Sleep(3 * time.Second)
	}

	time.Sleep(3 * time.Second)
}

func main() {
	vacuum.RegisterMain(Main)
	vacuum.RegisterString("MigrateTester", func() vacuum.StringDelegate {
		pt := &MigrateTester{}
		_ = vacuum.PersistentString(pt)
		return pt
	})
	vacuum_server.RunServer()
}
