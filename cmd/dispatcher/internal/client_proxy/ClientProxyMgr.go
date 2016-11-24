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
	sync.RWMutex // locker of the string ctrl

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

func registerClientProxyInfo(cp *ClientProxy, serverID int) (ret []int) {
	// register the new vacuum server client proxy
	cp.ServerID = serverID

	clientProxiesLock.Lock()
	clientProxes[serverID] = cp
	genClientProxyIDs()
	ret = clientProxyIDs
	vlog.Debug("registerClientProxyInfo: all client proxies: %v", clientProxes)
	clientProxiesLock.Unlock()
	return
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

func getStringCtrl(stringID string) (ret *StringCtrl) {
	stringCtrlsLock.RLock()
	ret = stringCtrls[stringID]
	if ret != nil { // common case
		stringCtrlsLock.RUnlock()
		return
	}
	// unlock and re-lock
	stringCtrlsLock.RUnlock()

	stringCtrlsLock.Lock()
	ret = stringCtrls[stringID]
	if ret == nil {
		ret = newStringCtrl()
		stringCtrls[stringID] = ret
	}
	stringCtrlsLock.Unlock()
	return
}

func setStringLocationMigrating(stringID string, serverID int, migrating bool) {
	ctrl := getStringCtrl(stringID)
	ctrl.Lock()
	ctrl.ServerID = serverID
	ctrl.Migrating = migrating
	ctrl.Unlock()

	//vlog.Debug("setStringLocationMigrating %s => %v", stringID, serverID)
}

func setStringLocation(stringID string, serverID int) {
	ctrl := getStringCtrl(stringID)
	ctrl.Lock()
	ctrl.ServerID = serverID
	ctrl.Unlock()

	//vlog.Debug("setStringLocation %s => %v", stringID, serverID)
}

func setStringMigrating(stringID string, migrating bool) {
	ctrl := getStringCtrl(stringID)
	ctrl.Lock()
	ctrl.Migrating = migrating
	ctrl.Unlock()

	//vlog.Debug("setStringMigrating %s => %v", stringID, migrating)
}
