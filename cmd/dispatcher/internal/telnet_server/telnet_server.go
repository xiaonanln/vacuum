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

//func serveTelnetServer() error {
//	defer func() {
//		if err := recover(); err != nil {
//			debuglog("panic: %v", err)
//		}
//	}()
//	ln, err := net.Listen("tcp", TELNET_SERVER_LISTEN_ATTR)
//	debuglog("Listening on %s ...", TELNET_SERVER_LISTEN_ATTR)
//
//	if err != nil {
//		return err
//	}
//
//	for {
//		conn, err := ln.Accept()
//		if err != nil {
//			if netutil.IsTemporaryNetError(err) {
//				continue
//			} else {
//				return err
//			}
//		}
//
//		go handleTelnetConnection(conn)
//	}
//
//	return nil
//}

func (d TelnetServerDelegate) ServeTCPConnection(conn net.Conn) {
	newTelnetConsole(conn).run()
}