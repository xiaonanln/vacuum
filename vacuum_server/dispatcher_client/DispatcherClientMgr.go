package dispatcher_client

import (
	"time"

	"sync/atomic"

	"unsafe"

	"errors"

	"strings"

	"strconv"

	"github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/config"
	"github.com/xiaonanln/vacuum/netutil"
	"github.com/xiaonanln/vacuum/vlog"
)

const (
	LOOP_DELAY_ON_DISPATCHER_CLIENT_ERROR = 3 * time.Second
)

var (
	_dispatcherClient         *DispatcherClient
	serverID                  = 0
	errDispatcherNotConnected = errors.New("dispatcher not connected")
)

func getDispatcherClient() *DispatcherClient {
	addr := (*uintptr)(unsafe.Pointer(&_dispatcherClient))
	return (*DispatcherClient)(unsafe.Pointer(atomic.LoadUintptr(addr)))
}

func setDispatcherClient(dc *DispatcherClient) {
	addr := (*uintptr)(unsafe.Pointer(&_dispatcherClient))
	atomic.StoreUintptr(addr, uintptr(unsafe.Pointer(dc)))
}

func assureConnectedDispatcherClient() *DispatcherClient {
	var err error
	dispatcherClient := getDispatcherClient()
	vlog.Debug("assureConnectedDispatcherClient: dispatcherClient", dispatcherClient)
	for dispatcherClient == nil {
		dispatcherClient, err = connectDispatchClient()
		if err != nil {
			vlog.Errorf("Connect to dispatcher failed: %s", err.Error())
			time.Sleep(LOOP_DELAY_ON_DISPATCHER_CLIENT_ERROR)
			continue
		}

		if serverID == 0 {
			vlog.Panicf("invalid serverID: %v", serverID)
		}

		dispatcherClient.RegisterVacuumServer(serverID)
		setDispatcherClient(dispatcherClient)
	}

	return dispatcherClient
}

func connectDispatchClient() (*DispatcherClient, error) {
	dispatcherPublicIP := config.GetConfig().Dispatcher.PublicIP
	port, err := strconv.Atoi(strings.Split(config.GetConfig().Dispatcher.Host, ":")[1])
	if err != nil {
		panic(err)
	}

	conn, err := netutil.ConnectTCP(dispatcherPublicIP, port)
	if err != nil {
		return nil, err
	}
	return newDispatcherClient(conn), nil
}

func Initialize(_serverID int, h DispatcherRespHandler) {
	serverID = _serverID
	dispatcherRespHandler = h

	assureConnectedDispatcherClient()
	go netutil.ServeForever(serveDispatcherClient)
}

func SendStringMessage(stringID string, msg common.StringMessage) error {
	return getDispatcherClientForSend().SendStringMessage(stringID, msg)
}

func SendCreateStringReq(name string, stringID string, args []interface{}) error {
	return getDispatcherClientForSend().SendCreateStringReq(name, stringID, args)
}

func SendLoadStringReq(name string, stringID string) error {
	return getDispatcherClientForSend().SendLoadStringReq(name, stringID)
}

func SendCreateStringLocallyReq(name string, stringID string) error {
	return getDispatcherClientForSend().SendCreateStringLocallyReq(name, stringID)
}

func SendDeclareServiceReq(stringID string, serviceName string) error {
	return getDispatcherClientForSend().SendDeclareServiceReq(stringID, serviceName)
}

func SendStringDelReq(stringID string) error {
	return getDispatcherClientForSend().SendStringDelReq(stringID)
}

func RelayCloseString(stringID string) error {
	return getDispatcherClientForSend().RelayCloseString(stringID)
}

func SendStartMigrateStringReq(stringID string) error {
	return getDispatcherClientForSend().SendStartMigrateStringReq(stringID)
}

func SendMigrateStringReq(name string, stringID string, serverID int, data map[string]interface{}) error {
	return getDispatcherClientForSend().SendMigrateStringReq(name, stringID, serverID, data)
}

func getDispatcherClientForSend() *DispatcherClient {
	dispatcherClient := getDispatcherClient()
	//if dispatcherClient == nil {
	//	debug.PrintStack()
	//	vlog.Errorf("dispatcher client is nil")
	//	return errDispatcherNotConnected
	//}
	return dispatcherClient
}

// serve the dispatcher client, receive RESPs from dispatcher and process
func serveDispatcherClient() {
	var err error
	vlog.Debugf("serveDispatcherClient: start serving dispatcher client ...")
	for {
		dispatcherClient := assureConnectedDispatcherClient()

		err = dispatcherClient.RecvMsg(dispatcherClient)
		if err != nil {
			vlog.Errorf("serveDispatcherClient: RecvMsgPacket error: %s", err.Error())
			dispatcherClient.Close()
			setDispatcherClient(nil)
			time.Sleep(LOOP_DELAY_ON_DISPATCHER_CLIENT_ERROR)
			continue
		}
	}
}
