package main

import (
	"time"

	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/ext/entity"
	"github.com/xiaonanln/vacuum/vacuum_server"
	"github.com/xiaonanln/vacuum/vlog"
)

type PingPongEntity struct {
	entity.BaseEntity // each entity class derives the BaseEntity
}

func (t *PingPongEntity) Ping(pongID entity.EntityID, v int) {
	vlog.Info("########################## %s.Ping %v ##########################", t, v)
	time.Sleep(500 * time.Millisecond)
	pongID.Call("Pong", t.ID, v+1)
}

func (t *PingPongEntity) Pong(pingID entity.EntityID, v int) {
	vlog.Info("########################## %s.Pong %v ##########################", t, v)
	time.Sleep(500 * time.Millisecond)
	pingID.Call("Ping", t.ID, v+1)
}

func Main(s *vacuum.String) {
	ping := entity.CreateEntity("PingPongEntity")
	pong := entity.CreateEntity("PingPongEntity")

	ping.Call("Ping", pong, 1)

	time.Sleep(10 * time.Second)
}

func main() {
	entity.RegisterEntity("PingPongEntity", &PingPongEntity{})
	vacuum.RegisterMain(Main)
	vacuum_server.RunServer()
}
