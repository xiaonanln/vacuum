package vacuum_server

import "time"

const (
	DISPATCHER_ADDR = ":"
)

var (
	dispatcherClient *DispatcherClient
)

func RunServer() {
	go manageDispatcherClient()

	dispatcherClient = newDispatcherClient()
	for {
		time.Sleep(time.Second)
	}
}
