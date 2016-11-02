package distributed_server

import (
	"net"

	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/netutil"
)

const (
	TCP_LISTENER_STRING_NAME = "client_conn._TCPListenerStringDelegateMaker"
	CONN_HANDLER_STRING_NAME = "client_conn._ConnHandlerStringDelegateMaker"
)

func init() {
	vacuum.RegisterString(TCP_LISTENER_STRING_NAME, _TCPListenerStringDelegateMaker)
	vacuum.RegisterString(CONN_HANDLER_STRING_NAME, _ConnHandlerStringDelegateMaker)
}

type TCPListener struct {
}

type ConnHandler func(conn net.Conn)

type _ServeTCPDelegate struct {
	handleConn ConnHandler
}

func _TCPListenerStringDelegateMaker(s *vacuum.String) {
	listenAddr := s.ReadString()
	handleConn := s.Read().(ConnHandler)
	netutil.ServeTCP(listenAddr, _ServeTCPDelegate{
		handleConn: handleConn,
	})
}

func StartServeTCP(addr string, handleConn ConnHandler) {
	listenerID := vacuum.CreateStringLocally(TCP_LISTENER_STRING_NAME)
	vacuum.Send(listenerID, addr)
	vacuum.Send(listenerID, handleConn)
}

func (d _ServeTCPDelegate) ServeTCPConnection(conn net.Conn) {
	// create a new string for serving the conn
	handlerID := vacuum.CreateStringLocally(CONN_HANDLER_STRING_NAME)
	vacuum.Send(handlerID, conn)
	vacuum.Send(handlerID, d.handleConn)
}

func _ConnHandlerStringDelegateMaker(s *vacuum.String) {
	conn := s.Read().(net.Conn)
	handleConn := s.Read().(ConnHandler)

	defer conn.Close()
	handleConn(conn)
}
