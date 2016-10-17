package client_proxy

import (
	"log"
	"math/rand"
	"sync"
)

var (
	clientProxyMgrLock = sync.RWMutex{}
	clientProxes       = map[int]*ClientProxy{}
	clientProxyIDs     = []int{}
)

func getRandomClientProxy() (ret *ClientProxy) {
	clientProxyMgrLock.RLock()
	i := rand.Intn(len(clientProxyIDs))
	ret = clientProxes[clientProxyIDs[i]]
	clientProxyMgrLock.RUnlock()
	return
}

func registerClientProxyInfo(cp *ClientProxy, serverID int) {
	// register the new vacuum server client proxy
	cp.ServerID = serverID

	clientProxyMgrLock.Lock()
	clientProxes[serverID] = cp
	genClientProxyIDs()
	clientProxyMgrLock.Unlock()

	log.Printf("registerClientProxyInfo: all client proxies: %v", clientProxes)
}

func onClientProxyClose(cp *ClientProxy) {
	// called when the vacuum server is disconnected
	serverID := cp.ServerID
	if serverID == 0 {
		// should not happen
		return
	}

	clientProxyMgrLock.Lock()

	if clientProxes[serverID] == cp {
		delete(clientProxes, serverID)
		genClientProxyIDs()
	}

	clientProxyMgrLock.Unlock()
	log.Printf("onClientProxyClose %v: all client proxies: %v", serverID, clientProxes)
}

func genClientProxyIDs() {
	clientProxyIDs = make([]int, 0, len(clientProxes))
	for id, _ := range clientProxes {
		clientProxyIDs = append(clientProxyIDs, id)
	}
}
