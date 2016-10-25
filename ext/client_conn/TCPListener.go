package client_conn

import (
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/netutil"
)

const (
	TCP_LISTENER_STRING_NAME = "client_conn._TCPListenerStringRoutine"
)

func init() {
	vacuum.RegisterString(TCP_LISTENER_STRING_NAME, _TCPListenerStringRoutine)
}

type TCPListener struct {
}

type _TCPListenerDelegate struct{}

func _TCPListenerStringRoutine(s *vacuum.String) {
	listenAddr := s.ReadString()
	netutil.ServeTCP(listenAddr, _TCPListenerDelegate{})
}

func NewTCPListener(addr string) {
	listenerID := vacuum.CreateStringLocally(TCP_LISTENER_STRING_NAME)
	vacuum.Send(listenerID, addr)
}

func (d _TCPListenerDelegate) ServeTCPConnection(conn net.Conn) {
	logrus.Debugf("New TCP connection: %s", conn.RemoteAddr())
}
