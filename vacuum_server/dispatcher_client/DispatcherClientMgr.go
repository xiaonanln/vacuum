package dispatcher_client

import (
	"log"

	"time"

	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/netutil"
)

var (
	dispatcherClient *DispatcherClient
)

func maintainDispatcherClient() {
	var err error

	for dispatcherClient == nil {
		dispatcherClient, err = connectDispatchClient()
		if err != nil {
			log.Printf("Connect to dispatcher failed: %s", err.Error())
			time.Sleep(time.Second)
		}
	}
}

func connectDispatchClient() (*DispatcherClient, error) {
	conn, err := netutil.ConnectTCP("localhost", 7581)
	if err != nil {
		return nil, err
	}
	return newDispatcherClient(conn), nil
}

func RegisterVacuumServer() {
	maintainDispatcherClient()
	dispatcherClient.RegisterVacuumServer()
}

func SendStringMessage(sid string, msg vacuum.StringMessage) {
	var err error
	maintainDispatcherClient()
	err = dispatcherClient.SendStringMessage(sid, msg)
	if err != nil {
		log.Printf("SendStringMessage: send string message failed with error %s, dispatcher lost ..", err.Error())
		dispatcherClient.Close()
		dispatcherClient = nil
	}
}
