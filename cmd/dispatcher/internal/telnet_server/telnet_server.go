package telnet_server

import (
	"fmt"
	"net"
	"sync"

	"github.com/xiaonanln/vacuum/config"
	"github.com/xiaonanln/vacuum/netutil"
	"github.com/xiaonanln/vacuum/vlog"
)

func debuglog(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	vlog.Debug("TelnetServer: %s", s)
}

type TelnetServerDelegate struct {
}

func ServeTelnetServer(wait *sync.WaitGroup) {
	netutil.ServeTCPForever(config.GetConfig().Dispatcher.ConsoleHost, &TelnetServerDelegate{})
	wait.Done()
}

func (d TelnetServerDelegate) ServeTCPConnection(conn net.Conn) {
	newTelnetConsole(conn).run()
}
