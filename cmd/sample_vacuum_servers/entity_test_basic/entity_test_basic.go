package main

import (
	"time"

	"github.com/xiaonanln/typeconv"
	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/ext/entity"
	"github.com/xiaonanln/vacuum/vacuum_server"
	"github.com/xiaonanln/vacuum/vlog"
)

type TestEntity struct {
	entity.Entity // each entity class derives the BaseEntity
	value         int
}

func (t *TestEntity) TestFunc1() {
	vlog.Info("#################################################################### %s.TestFunc1 %v ####################################################################", t, t.value)
	t.Save()

	t.Migrate(1)
}

func (t *TestEntity) TestFunc2(a int, b float64, c string) {
	vlog.Info("#################################################################### %s.TestFunc2 %v %v %v %v ####################################################################", t, t.value, a, b, c)
	t.Save()
}

func (t *TestEntity) OnMigrateOut(extra map[string]interface{}) {
	extra["value"] = t.value + 1
}

func (t *TestEntity) OnMigrateIn(extra map[string]interface{}) {
	t.value = int(typeconv.Int(extra["value"]))
}

func Main() {
	entityID := entity.CreateEntity("TestEntity")
	entityID.Call("TestFunc1")
	time.Sleep(1 * time.Second)
	entityID.Call("TestFunc2", 1, 2, "3")
	time.Sleep(1 * time.Second)
	//entityID.Call("TestFunc2", 1, 2)
	//time.Sleep(1 * time.Second)
	//
	//entityID.Call("TestFuncNotExist", 1, 2)
	//time.Sleep(1 * time.Second)

}

func main() {
	entity.RegisterEntity("TestEntity", &TestEntity{})
	vacuum.RegisterMain(Main)
	vacuum_server.RunServer()
}
