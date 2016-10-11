package telnet_server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/xiaonanln/vacuum/netutils"
)

func debug(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	log.Printf("TelnetServer: %s", s)
}

func ServeTelnetServer(wait *sync.WaitGroup) {
	for {
		err := serveTelnetServer()
		debug("error: %s", err.Error())
	}
	wait.Done()
}

func serveTelnetServer() error {
	defer func() {
		if err := recover(); err != nil {
			debug("panic: %v", err)
		}
	}()
	ln, err := net.Listen("tcp", ":7582")

	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			if netutils.IsTemporaryNetworkError(err) {
				continue
			} else {
				return err
			}
		}

		go handleTelnetConnection(conn)
	}

	return nil
}

func handleTelnetConnection(conn net.Conn) {
	newTelnetConsole(conn).run()
}
