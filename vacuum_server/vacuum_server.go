package vacuum_server

import (
	"time"

	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/common"
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

func (rh DispatcherRespHandler) HandleDispatcherResp_CreateString(name string, stringID string) {
	vacuum.OnCreateString(name, stringID)
}

func (rh DispatcherRespHandler) HandleDispatcherResp_DeclareService(stringID string, serviceName string) {
	vacuum.OnDeclareService(stringID, serviceName)
}

func (rh DispatcherRespHandler) HandleDispatcherResp_SendStringMessage(stringID string, msg common.StringMessage) {
	vacuum.OnSendStringMessage(stringID, msg)
}

func RunServer() {
	vacuum.CreateString("Main")
	for {
		time.Sleep(time.Second)
	}
}
