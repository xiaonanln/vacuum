package main

import (
	"sync"

	"fmt"

	"net"

	"flag"

	"github.com/xiaonanln/vacuum/cmd/dispatcher/internal/client_proxy"
	"github.com/xiaonanln/vacuum/cmd/dispatcher/internal/telnet_server"
	"github.com/xiaonanln/vacuum/config"
	"github.com/xiaonanln/vacuum/netutil"
	"github.com/xiaonanln/vacuum/vlog"
)

var (
	configFile = ""
)

func debuglog(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	vlog.Debugf("dispatcher: %s", s)
}

type DispatcherDelegate struct{}

func main() {
	flag.StringVar(&configFile, "c", config.CONFIG_FILENAME, "config file")
	flag.Parse()

	config.LoadConfig(configFile)
	vlog.SetLevel(vlog.DEBUG)

	wait := &sync.WaitGroup{}
	wait.Add(1)
	go telnet_server.ServeTelnetServer(wait)

	netutil.ServeTCPForever(config.GetConfig().Dispatcher.Host, &DispatcherDelegate{})
	wait.Wait()
}

func (dd *DispatcherDelegate) ServeTCPConnection(conn net.Conn) {
	client_proxy.NewClientProxy(conn).Serve()
}
