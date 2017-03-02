package vacuum_server

import (
	"time"

	"flag"

	"os"

	"sync"

	"github.com/xiaonanln/goTimer"
	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/config"
	"github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"
	"github.com/xiaonanln/vacuum/vlog"
)

const (
	DISPATCHER_ADDR   = ":"
	DEFAULT_LOG_LEVEL = "info"
)

var (
	serverID   = 1 // default server ID to be 1
	logLevel   = DEFAULT_LOG_LEVEL
	configFile = ""
	serverOps  _ServerOps
)

func init() {
	// initializing the vacuum server
	flag.IntVar(&serverID, "sid", 1, "server ID")
	flag.StringVar(&logLevel, "log", DEFAULT_LOG_LEVEL, "log level")
	flag.StringVar(&configFile, "c", config.DEFAULT_CONFIG_FILENAME, "config file")
	flag.Parse()

	setupLog(logLevel)

	if serverID <= 0 {
		vlog.Panicf("Server ID must be positive, not %d", serverID)
	}

	vlog.Debug(">>> Server ID: %d, Config file: %s", serverID, configFile)

	config.LoadConfig(configFile)
	vacuumConfig := config.GetConfig().GetVacuumConfig(serverID)
	vlog.Debug("VACUUM %d LOAD CONFIG:", serverID)
	os.Stderr.WriteString(config.FormatConfig(vacuumConfig))

	timer.StartTicks(10 * time.Millisecond)

	storage := openStorage(vacuumConfig.Storage)
	vacuum.Setup(serverID, storage, &serverOps)
}

type DispatcherRespHandler struct{}

func (rh DispatcherRespHandler) HandleDispatcherResp_RegisterVacuumServer(serverIDs []int) {
	vlog.Info("RegisterVacuumServer: %v", serverIDs)
	serverOps.setVacuumServerIDs(serverIDs)
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

func (rh DispatcherRespHandler) HandleDispatcherResp_LoadString(name string, stringID string) {
	vacuum.OnLoadString(name, stringID)
}

func (rh DispatcherRespHandler) HandleDispatcherResp_MigrateString(stringID string) {
	vacuum.MigrateString(stringID)
}

func (rh DispatcherRespHandler) HandleDispatcherResp_OnMigrateString(name string, stringID string, initArgs []interface{}, data map[string]interface{}, extraMigrateInfo map[string]interface{}) {
	vacuum.OnMigrateString(name, stringID, initArgs, data, extraMigrateInfo)
}

type _ServerOps struct {
	lock      sync.RWMutex
	serverIDs []int
}

func (so *_ServerOps) GetServerList() (ret []int) {
	so.lock.RLock()
	//vlog.Debug("GetServerList: %v", so.serverIDs)
	ret = so.serverIDs
	so.lock.RUnlock()
	return
}

func (so *_ServerOps) setVacuumServerIDs(serverIDs []int) {
	so.lock.Lock()
	so.serverIDs = serverIDs
	so.lock.Unlock()
}

func RunServer() {
	dispatcher_client.Initialize(serverID, DispatcherRespHandler{})

	vacuum.CreateStringLocally("Main")

	for {
		time.Sleep(time.Second)
	}
}

func ServerID() int {
	return serverID
}
