package netutil

import (
	"net"

	"log"
	"time"
)

const (
	RESTART_TCP_SERVER_INTERVAL = 3 * time.Second
)

type TCPServerDelegate interface {
	ServeTCPConnection(net.Conn)
}

func ServeTCP(listenAddr string, delegate TCPServerDelegate) {
	for {
		err := serveTCPImpl(listenAddr, delegate)
		log.Printf("server@%s failed with error: %v, will restart after %s", listenAddr, err, RESTART_TCP_SERVER_INTERVAL)
		time.Sleep(RESTART_TCP_SERVER_INTERVAL)
	}
}

func serveTCPImpl(listenAddr string, delegate TCPServerDelegate) error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic: %v", err)
		}
	}()

	ln, err := net.Listen("tcp", listenAddr)
	log.Printf("Listening on %s ...", listenAddr)

	if err != nil {
		return err
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if IsTemporaryNetError(err) {
				continue
			} else {
				return err
			}
		}

		log.Printf("Connection from: %s", conn.RemoteAddr())
		go delegate.ServeTCPConnection(conn)
	}
	return nil
}
