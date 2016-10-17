package telnet_server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/xiaonanln/vacuum/netutil"
)

const (
	TELNET_SERVER_LISTEN_ATTR = ":7582"
)

func debuglog(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	log.Printf("TelnetServer: %s", s)
}

type TelnetServerDelegate struct {
}

func ServeTelnetServer(wait *sync.WaitGroup) {
	netutil.ServeTCP(TELNET_SERVER_LISTEN_ATTR, &TelnetServerDelegate{})
	wait.Done()
}

func (d TelnetServerDelegate) ServeTCPConnection(conn net.Conn) {
	newTelnetConsole(conn).run()
}
