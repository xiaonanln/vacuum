package telnet_server

import "net"

type TelnetConsole struct {
	net.Conn
}

func newTelnetConsole(conn net.Conn) *TelnetConsole {
	tc := &TelnetConsole{
		Conn: conn,
	}
	return tc
}

func (tc *TelnetConsole) run() {
	tc.Close()
}
