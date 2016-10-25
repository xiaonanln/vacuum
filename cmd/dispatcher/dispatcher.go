package main

import (
	"sync"

	"fmt"

	log "github.com/Sirupsen/logrus"

	"net"

	"github.com/xiaonanln/vacuum/cmd/dispatcher/internal/client_proxy"
	"github.com/xiaonanln/vacuum/cmd/dispatcher/internal/telnet_server"
	"github.com/xiaonanln/vacuum/config"
	"github.com/xiaonanln/vacuum/netutil"
)

const (
	DISPATCHER_SERVER_LISTEN_ATTR = ":7581"
)

func debuglog(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	log.Debugf("dispatcher: %s", s)
}

type DispatcherDelegate struct{}

func main() {
	config.LoadConfig()
	log.SetLevel(log.DebugLevel)

	wait := &sync.WaitGroup{}
	wait.Add(1)
	go telnet_server.ServeTelnetServer(wait)

	netutil.ServeTCPForever(DISPATCHER_SERVER_LISTEN_ATTR, &DispatcherDelegate{})
	wait.Wait()
}

func (dd *DispatcherDelegate) ServeTCPConnection(conn net.Conn) {
	client_proxy.NewClientProxy(conn).Serve()
}
