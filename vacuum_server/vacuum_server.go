package vacuum_server

import (
	"time"

	"flag"

	log "github.com/Sirupsen/logrus"

	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/config"
	"github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"
)

const (
	DISPATCHER_ADDR = ":"
)

var (
	serverID = 1 // default server ID to be 1
)

type DispatcherRespHandler struct{}

func init() {
	// initializing the vacuum server
	flag.IntVar(&serverID, "sid", 1, "server ID")
	flag.Parse()

	setupLog()

	if serverID <= 0 {
		log.Panicf("Server ID must be positive, not %d", serverID)
	}

	log.Printf(">>> Server ID: %d", serverID)

	config.LoadConfig()
	dispatcher_client.Initialize(serverID, DispatcherRespHandler{})
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
	vacuum.CreateStringLocally("Main")
	for {
		time.Sleep(time.Second)
	}
}

func ServerID() int {
	return serverID
}
