package netutil

import (
	"net"

	"runtime/debug"
	"time"

	log "github.com/Sirupsen/logrus"
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
		log.Errorf("server@%s failed with error: %v, will restart after %s", listenAddr, err, RESTART_TCP_SERVER_INTERVAL)
		time.Sleep(RESTART_TCP_SERVER_INTERVAL)
	}
}

func serveTCPImpl(listenAddr string, delegate TCPServerDelegate) error {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("serveTCPImpl: paniced with error %s", err)
			debug.PrintStack()
		}
	}()

	ln, err := net.Listen("tcp", listenAddr)
	log.WithFields(log.Fields{"addr": listenAddr}).Info("Listening on TCP ...")

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

		log.Infof("Connection from: %s", conn.RemoteAddr())
		go delegate.ServeTCPConnection(conn)
	}
	return nil
}
