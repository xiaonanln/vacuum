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
	vlog.Info("TestFunc1")
}

func Main(s *vacuum.String) {
	entityID := entity.CreateEntity("TestEntity")
	entityID.Call("TestFunc1")
	time.Sleep(3 * time.Second)
}

func main() {
	entity.RegisterEntity("TestEntity", TestEntity{})
	vacuum.RegisterMain(Main)
	vacuum_server.RunServer()
}
