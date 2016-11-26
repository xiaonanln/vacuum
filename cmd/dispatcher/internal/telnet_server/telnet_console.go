package telnet_server

import (
	"net"
	"sync"

	"strings"

	"runtime/debug"

	"runtime"

	"fmt"

	"github.com/xiaonanln/vacuum/netutil"
	"github.com/xiaonanln/vacuum/proto"
	"github.com/xiaonanln/vacuum/vlog"
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
	vlog.Info("Console %s closed.", tc.conn.RemoteAddr())
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
			vlog.Debug("TelnetConsole.handleCommand failed: cmd=%s, err=%s", cmd, err)
			debug.PrintStack()
		}
	}()

	if cmd == "quit" || cmd == "exit" {
		tc.handleQuit()
	} else if cmd == "gc" {
		tc.handleGC()
	} else {
		tc.handleUnknownCommand(cmd)
	}
}

func (tc *TelnetConsole) handleQuit() {
	tc.close()
}

var (
	syncPool = sync.Pool{
		New: func() interface{} {
			return &proto.Message{}
		},
	}
)

func (tc *TelnetConsole) handleGC() {
	var messages []*proto.Message
	for i := 0; i < 1000; i++ {
		m := syncPool.Get().(*proto.Message)
		messages = append(messages, m)
	}
	for _, m := range messages {
		syncPool.Put(m)
	}
	messages = nil
	runtime.GC() // get up-to-date statistics
}

func (tc *TelnetConsole) handleUnknownCommand(cmd string) {
	tc.writeLine(fmt.Sprintf("unknown command: %s", cmd))
}
