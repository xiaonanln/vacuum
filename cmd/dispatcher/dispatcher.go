package main

import (
	"sync"

	"github.com/xiaonanln/vacuum/cmd/dispatcher/telnet_server"
)

func main() {
	wait := &sync.WaitGroup{}
	wait.Add(1)
	telnet_server.ServeTelnetServer(wait)
}
