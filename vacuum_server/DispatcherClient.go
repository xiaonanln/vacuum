package vacuum_server

import "github.com/xiaonanln/vacuum/netutil"

type DispatcherClient struct {
	netutil.BinaryConnection
}

func newDispatcherClient() *DispatcherClient {
	return &DispatcherClient{}
}

func (dc *DispatcherClient) connect() {

}
