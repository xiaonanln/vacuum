package main

import (
	"sync"

	"fmt"
	"log"

	"github.com/xiaonanln/vacuum/cmd/dispatcher/internal/telnet_server"
)

const (
	DISPATCHER_SERVER_LISTEN_ATTR = ":7581"
)

func debuglog(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	log.Printf("dispatcher: %s", s)
}

func main() {
	wait := &sync.WaitGroup{}
	wait.Add(1)
	go telnet_server.ServeTelnetServer(wait)

	//serveDispatcher()
	wait.Wait()
}

//func serveDispatcher() {
//	defer func() {
//		if err := recover(); err != nil {
//			debuglog("panic: %v", err)
//		}
//	}()
//
//	ln, err := net.Listen("tcp", DISPATCHER_SERVER_LISTEN_ATTR)
//	debuglog("Listening on %s ...", DISPATCHER_SERVER_LISTEN_ATTR)
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
//	return nil
//}
