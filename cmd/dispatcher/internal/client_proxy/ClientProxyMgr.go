package client_proxy

import (
	"math/rand"
	"sync"

	"github.com/xiaonanln/vacuum/proto"
	"github.com/xiaonanln/vacuum/vlog"
)

type _CachedMessage struct {
	msg     *proto.Message
	pktsize uint32
}

type StringCtrl struct {
	ServerID       int  // Server ID of string
	Migrating      bool // if the string is migrating
	cachedMessages []_CachedMessage
}

func newStringCtrl() *StringCtrl {
	return &StringCtrl{
		cachedMessages: []_CachedMessage{},
	}
}

var (
	// ServerID to ClientProxy
	clientProxiesLock = sync.RWMutex{}
	clientProxes      = map[int]*ClientProxy{}
	clientProxyIDs    = []int{}

	// StringID to ServerID
	stringCtrlsLock = sync.RWMutex{}
	stringCtrls     = map[string]*StringCtrl{}
)

func getRandomClientProxy() (ret *ClientProxy) {
	clientProxiesLock.RLock()
	i := rand.Intn(len(clientProxyIDs))
	ret = clientProxes[clientProxyIDs[i]]
	clientProxiesLock.RUnlock()
	return
}

func getClientProxy(serverID int) (ret *ClientProxy) {
	clientProxiesLock.RLock()
	ret = clientProxes[serverID]
	//vlog.Debug("getClientProxy", clientProxes, serverID, "=>", ret)
	clientProxiesLock.RUnlock()
	return
}

func registerClientProxyInfo(cp *ClientProxy, serverID int) {
	// register the new vacuum server client proxy
	cp.ServerID = serverID

	clientProxiesLock.Lock()
	clientProxes[serverID] = cp
	genClientProxyIDs()
	vlog.Debug("registerClientProxyInfo: all client proxies: %v", clientProxes)
	clientProxiesLock.Unlock()
}

func onClientProxyClose(cp *ClientProxy) {
	// called when the vacuum server is disconnected
	serverID := cp.ServerID
	if serverID == 0 {
		// should not happen
		return
	}

	clientProxiesLock.Lock()

	if clientProxes[serverID] == cp {
		delete(clientProxes, serverID)
		genClientProxyIDs()
	}

	vlog.Debug("onClientProxyClose %v: all client proxies: %v", serverID, clientProxes)
	clientProxiesLock.Unlock()
}

func genClientProxyIDs() {
	clientProxyIDs = make([]int, 0, len(clientProxes))
	for id := range clientProxes {
		clientProxyIDs = append(clientProxyIDs, id)
	}
}

func lockStringCtrlForRead(stringID string) (ret *StringCtrl) {
	stringCtrlsLock.RLock()
	ret = stringCtrls[stringID]
	if ret != nil {
		return
	}

	stringCtrlsLock.RUnlock()
	stringCtrlsLock.Lock()
	ret = stringCtrls[stringID]
	if ret == nil {
		ret = newStringCtrl()
		stringCtrls[stringID] = ret
	}
	stringCtrlsLock.Unlock()

	stringCtrlsLock.RLock() // assure RLock
	return
}

func lockStringCtrlForWrite(stringID string) (ret *StringCtrl) {
	stringCtrlsLock.Lock()
	ret = stringCtrls[stringID]
	if ret == nil {
		ret = newStringCtrl()
		stringCtrls[stringID] = ret
	}
	// lock left locked
	return
}

func setStringLocationMigrating(stringID string, serverID int, migrating bool) {
	ctrl := lockStringCtrlForWrite(stringID)
	ctrl.ServerID = serverID
	ctrl.Migrating = migrating
	stringCtrlsLock.Unlock()

	//vlog.Debug("setStringLocationMigrating %s => %v", stringID, serverID)
}

func setStringLocation(stringID string, serverID int) {
	ctrl := lockStringCtrlForWrite(stringID)
	ctrl.ServerID = serverID
	stringCtrlsLock.Unlock()

	//vlog.Debug("setStringLocation %s => %v", stringID, serverID)
}

func setStringMigrating(stringID string, migrating bool) {
	ctrl := lockStringCtrlForWrite(stringID)
	ctrl.Migrating = migrating
	stringCtrlsLock.Unlock()

	//vlog.Debug("setStringMigrating %s => %v", stringID, migrating)
}
