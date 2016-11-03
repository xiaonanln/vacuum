package vacuum_server

import (
	"time"

	"flag"

	log "github.com/Sirupsen/logrus"

	"os"

	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/config"
	"github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"
)

const (
	DISPATCHER_ADDR   = ":"
	DEFAULT_LOG_LEVEL = "info"
)

var (
	serverID   = 1 // default server ID to be 1
	logLevel   = DEFAULT_LOG_LEVEL
	configFile = ""
)

type DispatcherRespHandler struct{}

func init() {
	// initializing the vacuum server
	flag.IntVar(&serverID, "sid", 1, "server ID")
	flag.StringVar(&logLevel, "log", DEFAULT_LOG_LEVEL, "log level")
	flag.StringVar(&configFile, "c", config.CONFIG_FILENAME, "config file")
	flag.Parse()

	setupLog(logLevel)

	if serverID <= 0 {
		log.Panicf("Server ID must be positive, not %d", serverID)
	}

	log.Debugf(">>> Server ID: %d, Config file: %s", serverID, configFile)

	config.LoadConfig(configFile)
	vacuumConfig := config.GetConfig().GetVacuumConfig(serverID)
	log.Printf("VACUUM %d LOAD CONFIG:", serverID)
	os.Stderr.WriteString(config.FormatConfig(vacuumConfig))

	storage := openStorage(vacuumConfig.Storage)
	vacuum.Setup(serverID, storage)

	dispatcher_client.Initialize(serverID, DispatcherRespHandler{})
}

func (rh DispatcherRespHandler) HandleDispatcherResp_CreateString(name string, stringID string, args []interface{}) {
	vacuum.OnCreateString(name, stringID, args)
}

func (rh DispatcherRespHandler) HandleDispatcherResp_DeclareService(stringID string, serviceName string) {
	vacuum.OnDeclareService(stringID, serviceName)
}

func (rh DispatcherRespHandler) HandleDispatcherResp_SendStringMessage(stringID string, msg common.StringMessage) {
	vacuum.OnSendStringMessage(stringID, msg)
}

func (rh DispatcherRespHandler) HandleDispatcherResp_CloseString(stringID string) {
	vacuum.OnCloseString(stringID)
}

func (rh DispatcherRespHandler) HandleDispatcherResp_DelString(stringID string) {
	vacuum.OnDelString(stringID)
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
