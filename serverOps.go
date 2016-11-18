package vacuum

import (
	"time"

	"github.com/xiaonanln/vacuum/vlog"
)

var (
	serverOps ServerOps
)

type ServerOps interface {
	GetServerList() []int
}

func WaitServerReady(serverID int) {
	infoed := false
	for {
		for _, _serverID := range serverOps.GetServerList() {
			//vlog.Info("%v %v", _serverID, serverID)
			if _serverID == serverID {
				return // server is ready
			}
		}
		// server not found, wait
		if !infoed {
			vlog.Info("Server %d is not ready, waiting ...", serverID)
			infoed = true
		}
		time.Sleep(WAIT_OPS_LOOP_INTERVAL)
	}
}
