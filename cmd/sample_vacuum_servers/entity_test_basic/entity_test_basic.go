package main

import (
	"time"

	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/ext/entity"
	"github.com/xiaonanln/vacuum/vacuum_server"
	"github.com/xiaonanln/vacuum/vlog"
)

type TestEntity struct {
	entity.BaseEntity // each entity class derives the BaseEntity
}

func (t *TestEntity) TestFunc1() {
	vlog.Info("#################################################################### %s.TestFunc1 ####################################################################", t)
}

func (t *TestEntity) TestFunc2(a int, b float64, c string) {
	vlog.Info("#################################################################### %s.TestFunc2 %v %v %v ####################################################################", t, a, b, c)
}

func Main(s *vacuum.String) {
	entityID := entity.CreateEntity("TestEntity")
	entityID.Call("TestFunc1")
	time.Sleep(1 * time.Second)
	entityID.Call("TestFunc2", 1, 2, "3")
	time.Sleep(1 * time.Second)
	entityID.Call("TestFunc2", 1, 2)
	time.Sleep(1 * time.Second)

	entityID.Call("TestFuncNotExist", 1, 2)
	time.Sleep(1 * time.Second)
}

func main() {
	entity.RegisterEntity("TestEntity", TestEntity{})
	vacuum.RegisterMain(Main)
	vacuum_server.RunServer()
}
