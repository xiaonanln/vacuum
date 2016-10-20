package telnet_server

import (
	"net"

	"strings"

	"fmt"

	log "github.com/Sirupsen/logrus"

	"runtime/debug"

	"github.com/xiaonanln/vacuum/netutil"
)

type TelnetConsole struct {
	conn net.Conn
}

func newTelnetConsole(conn net.Conn) *TelnetConsole {
	tc := &TelnetConsole{
		conn: conn,
	}
	tc.conn.(*net.TCPConn).SetNoDelay(true)
	return tc
}

func (tc *TelnetConsole) run() {
	defer tc.close()

	tc.writeLine("Welcome to vacuum!")
	for {
		line, err := tc.readLine()
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tc.handleCommand(line)
	}
	log.Printf("Console %s closed.", tc.conn.RemoteAddr())
}

func (tc *TelnetConsole) close() {
	tc.conn.Close()
}

func (tc *TelnetConsole) writeStr(s string) {
	netutil.WriteAll(tc.conn, []byte(s))
}

func (tc *TelnetConsole) writeLine(s string) {
	tc.writeStr(s + netutil.NEWLINE)
}

func (tc *TelnetConsole) readLine() (string, error) {
	tc.writeStr(">>> ")
	line, err := netutil.ReadLine(tc.conn)
	return line, err
}

func (tc *TelnetConsole) handleCommand(cmd string) {
	defer func() {
		err := recover() // catch all errors during handling command
		if err != nil {
			log.Printf("TelnetConsole.handleCommand failed: cmd=%s, err=%s", cmd, err)
			debug.PrintStack()
		}
	}()

	if cmd == "quit" || cmd == "exit" {
		tc.handleQuit()
	} else {
		tc.handleUnknownCommand(cmd)
	}
}

func (tc *TelnetConsole) handleQuit() {
	tc.close()
}

func (tc *TelnetConsole) handleUnknownCommand(cmd string) {
	tc.writeLine(fmt.Sprintf("unknown command: %s", cmd))
}
