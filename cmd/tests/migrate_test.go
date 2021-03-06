package test

import (
	"testing"

	"time"

	"math/rand"

	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/vacuum_server"
	"github.com/xiaonanln/vacuum/vlog"
)

func init() {
	vacuum.RegisterString("Main", Main)
	vacuum.RegisterString("TestString", _TestString)
}

func Main(s *vacuum.String) {
	vlog.Printf("Main running...")
}

type _TestStorage struct {
	MagicNum int64
}

func (ts *_TestStorage) GetMigrateData() interface{} {
	return map[string]interface{}{
		"MagicNum": ts.MagicNum,
	}
}

func (ts *_TestStorage) InitWithMigrateData(data interface{}) {
	m := data.(map[string]interface{})
	ts.MagicNum = m["MagicNum"].(int64)
}

func _TestString(s *vacuum.String) {
	ts := &_TestStorage{}
	ts.MagicNum = rand.Int63()

	vlog.Printf(">>> _TestString running, the magic number is %v <<<", ts.MagicNum)

	vlog.Printf(">>> Migrate after seconds ... <<<")
	time.Sleep(time.Second)

	s.Migrate(1)
}

func TestMigrate(t *testing.T) {
	vacuum.CreateString("TestString")
	vacuum_server.RunServer()
}
