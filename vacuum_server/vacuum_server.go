package vacuum_server

import (
	"sync"

	"github.com/xiaonanln/vacuum/cmd/dispatcher/telnet_server"
)

func RunServer() {
	wait := &sync.WaitGroup{}
	wait.Add(1)

	go telnet_server.ServeTelnetServer(wait) // new goroutine for telnet server
	wait.Wait()
}
