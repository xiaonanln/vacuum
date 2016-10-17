package dispatcher_client

import (
	"log"

	"time"

	"sync/atomic"

	"unsafe"

	"github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/netutil"
)

var (
	_dispatcherClient *DispatcherClient
	serverID          = 0
)

func getDispatcherClient() *DispatcherClient {
	addr := (*uintptr)(unsafe.Pointer(&_dispatcherClient))
	return (*DispatcherClient)(unsafe.Pointer(atomic.LoadUintptr(addr)))
}

func setDispatcherClient(dc *DispatcherClient) {
	addr := (*uintptr)(unsafe.Pointer(&_dispatcherClient))
	atomic.StoreUintptr(addr, uintptr(unsafe.Pointer(dc)))
}

func maintainDispatcherClient() *DispatcherClient {
	var err error
	dispatcherClient := getDispatcherClient()
	log.Println("dispatcherClient", dispatcherClient)
	for dispatcherClient == nil {
		dispatcherClient, err = connectDispatchClient()
		if err != nil {
			log.Printf("Connect to dispatcher failed: %s", err.Error())
			time.Sleep(time.Second)
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

func RegisterVacuumServer(_serverID int) {
	serverID = _serverID
	maintainDispatcherClient()
}

func SendStringMessage(sid string, msg common.StringMessage) {
	var err error
	dispatcherClient := maintainDispatcherClient()
	err = dispatcherClient.SendStringMessage(sid, msg)
	if err != nil {
		log.Printf("SendStringMessage: send string message failed with error %s, dispatcher lost ..", err.Error())
		dispatcherClient.Close()
		setDispatcherClient(nil)
	}
}

func CreateString(name string) error {
	dispatcherClient := maintainDispatcherClient()
	return dispatcherClient.CreateString(name)
}
