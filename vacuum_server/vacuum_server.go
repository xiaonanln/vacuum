package vacuum_server

import (
	"time"

	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/config"
	"github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"
)

const (
	DISPATCHER_ADDR = ":"
)

type DispatcherRespHandler struct{}

func init() {
	// initializing the vacuum server
	config.LoadConfig()
	dispatcher_client.Initialize(1, DispatcherRespHandler{})
}

func (rh DispatcherRespHandler) HandleDispatcherResp_CreateString(name string) {
	vacuum.OnCreateString(name)
}

func (rh DispatcherRespHandler) HandleDispatcherResp_DeclareService(stringID string, serviceName string) {
	vacuum.OnDeclareService(stringID, serviceName)
}

func RunServer() {
	for {
		time.Sleep(time.Second)
	}
}
