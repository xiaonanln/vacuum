package main

import (
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/ext/distributed_server"
	"github.com/xiaonanln/vacuum/netutil"
	"github.com/xiaonanln/vacuum/vacuum_server"
)

func isPrimaryServer() bool {
	return vacuum_server.ServerID() == 1
}

func Main(s *vacuum.String) {
	distributed_server.StartServeTCP(":12345", handleConn)
}

func checkConnError(err error) {
	if err != nil {
		panic(err)
	}
}

func handleConn(conn net.Conn) {
	logrus.Printf("New client connection: %s", conn.RemoteAddr())
	err := netutil.WriteAll(conn, []byte("Welcome to echo server\n"))
	checkConnError(err)

	var data [1024]byte
	for {
		n, err := conn.Read(data[:])
		checkConnError(err)
		netutil.WriteAll(conn, data[:n])
	}
}

func main() {
	vacuum.RegisterString("Main", Main)
	vacuum_server.RunServer()
}
