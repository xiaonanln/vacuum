package dispatcher_client

import (
	"log"

	"time"

	"github.com/xiaonanln/vacuum/common"
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

func RegisterVacuumServer(serverID int) {
	maintainDispatcherClient()
	dispatcherClient.RegisterVacuumServer(serverID)
}

func SendStringMessage(sid string, msg common.StringMessage) {
	maintainDispatcherClient()

	var err error
	err = dispatcherClient.SendStringMessage(sid, msg)
	if err != nil {
		log.Printf("SendStringMessage: send string message failed with error %s, dispatcher lost ..", err.Error())
		dispatcherClient.Close()
		dispatcherClient = nil
	}
}

func CreateString(name string) error {
	maintainDispatcherClient()
	return dispatcherClient.CreateString(name)
}
