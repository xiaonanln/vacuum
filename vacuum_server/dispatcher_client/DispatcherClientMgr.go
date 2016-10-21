package dispatcher_client

import (
	log "github.com/Sirupsen/logrus"

	"time"

	"sync/atomic"

	"unsafe"

	"errors"

	"runtime/debug"

	"github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/netutil"
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
	log.Debugln("assureConnectedDispatcherClient: dispatcherClient", dispatcherClient)
	for dispatcherClient == nil {
		dispatcherClient, err = connectDispatchClient()
		if err != nil {
			log.Errorf("Connect to dispatcher failed: %s", err.Error())
			time.Sleep(LOOP_DELAY_ON_DISPATCHER_CLIENT_ERROR)
			continue
		}

		if serverID == 0 {
			log.Panicf("invalid serverID: %v", serverID)
		}

		dispatcherClient.RegisterVacuumServer(serverID)
		setDispatcherClient(dispatcherClient)
	}

	return dispatcherClient
}

func connectDispatchClient() (*DispatcherClient, error) {
	conn, err := netutil.ConnectTCP("localhost", 7581)
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
	dispatcherClient := getDispatcherClient()
	if dispatcherClient == nil {
		debug.PrintStack()
		log.Errorf("dispatcher client is nil")
		return errDispatcherNotConnected
	}
	return dispatcherClient.SendStringMessage(stringID, msg)
}

func SendCreateStringReq(name string, stringID string) error {
	dispatcherClient := getDispatcherClient()
	if dispatcherClient == nil {
		debug.PrintStack()
		log.Errorf("dispatcher client is nil")
		return errDispatcherNotConnected
	}
	return dispatcherClient.SendCreateStringReq(name, stringID)
}

func SendCreateStringLocallyReq(name string, stringID string) error {
	dispatcherClient := getDispatcherClient()
	if dispatcherClient == nil {
		debug.PrintStack()
		log.Errorf("dispatcher client is nil")
		return errDispatcherNotConnected
	}
	return dispatcherClient.SendCreateStringLocallyReq(name, stringID)
}

func SendDeclareServiceReq(stringID string, serviceName string) error {
	dispatcherClient := getDispatcherClient()
	if dispatcherClient == nil {
		debug.PrintStack()
		log.Errorf("dispatcher client is nil")
		return errDispatcherNotConnected
	}
	return dispatcherClient.SendDeclareServiceReq(stringID, serviceName)
}

func SendCloseStringReq(stringID string) error {
	dispatcherClient := getDispatcherClient()
	if dispatcherClient == nil {
		debug.PrintStack()
		log.Errorf("dispatcher client is nil")
		return errDispatcherNotConnected
	}
	return dispatcherClient.SendCloseStringReq(stringID)
}

// serve the dispatcher client, receive RESPs from dispatcher and process
func serveDispatcherClient() {
	var err error
	log.Debugf("serveDispatcherClient: start serving dispatcher client ...")
	for {
		dispatcherClient := assureConnectedDispatcherClient()

		err = dispatcherClient.RecvMsg(dispatcherClient)
		if err != nil {
			log.Errorf("serveDispatcherClient: RecvMsgPacket error: %s", err.Error())
			dispatcherClient.Close()
			setDispatcherClient(nil)
			time.Sleep(LOOP_DELAY_ON_DISPATCHER_CLIENT_ERROR)
			continue
		}
	}
}
