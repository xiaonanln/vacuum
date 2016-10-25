package main

import (
	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/ext/client_conn"
	"github.com/xiaonanln/vacuum/vacuum_server"
)

func isPrimaryServer() bool {
	return vacuum_server.ServerID() == 1
}

func Main(s *vacuum.String) {
	client_conn.NewTCPListener(":12345")
}

func main() {
	vacuum.RegisterString("Main", Main)
	vacuum_server.RunServer()
}
