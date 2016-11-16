package client_proxy

import (
	"math/rand"
	"sync"

	"github.com/xiaonanln/vacuum/vlog"
)

type _StringInfo struct {
	ServerID  int
	Migrating bool
}

var (
	// ServerID to ClientProxy
	clientProxiesLock = sync.RWMutex{}
	clientProxes      = map[int]*ClientProxy{}
	clientProxyIDs    = []int{}

	// StringID to ServerID
	stringInfosLock = sync.RWMutex{}
	stringInfos     = map[string]_StringInfo{}
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
	vlog.Debugf("registerClientProxyInfo: all client proxies: %v", clientProxes)
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

	vlog.Debugf("onClientProxyClose %v: all client proxies: %v", serverID, clientProxes)
	clientProxiesLock.Unlock()
}

func genClientProxyIDs() {
	clientProxyIDs = make([]int, 0, len(clientProxes))
	for id := range clientProxes {
		clientProxyIDs = append(clientProxyIDs, id)
	}
}

func setStringLocationMigrating(stringID string, serverID int, migrating bool) {
	stringInfosLock.Lock()
	info := stringInfos[stringID]
	info.ServerID = serverID
	info.Migrating = migrating
	stringInfos[stringID] = info
	stringInfosLock.Unlock()

	//vlog.Debugf("setStringLocationMigrating %s => %v", stringID, serverID)
}

func setStringLocation(stringID string, serverID int) {
	stringInfosLock.Lock()
	info := stringInfos[stringID]
	info.ServerID = serverID
	stringInfos[stringID] = info
	stringInfosLock.Unlock()

	//vlog.Debugf("setStringLocation %s => %v", stringID, serverID)
}

func setStringMigrating(stringID string, migrating bool) {
	stringInfosLock.Lock()
	info := stringInfos[stringID]
	info.Migrating = migrating
	stringInfos[stringID] = info
	stringInfosLock.Unlock()

	//vlog.Debugf("setStringMigrating %s => %v", stringID, migrating)
}

func getStringInfo(stringID string) (ret _StringInfo) {
	stringInfosLock.RLock()
	ret = stringInfos[stringID]
	stringInfosLock.RUnlock()

	//vlog.Debugf("getStringInfo %s => %v", stringID, ret)
	return
}

func getStringLocation(stringID string) int {
	return getStringInfo(stringID).ServerID
}
