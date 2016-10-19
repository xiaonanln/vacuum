package client_proxy

import (
	"log"
	"math/rand"
	"sync"
)

var (
	clientProxiesLock = sync.RWMutex{}
	clientProxes      = map[int]*ClientProxy{}
	clientProxyIDs    = []int{}

	stringLocationsLock = sync.RWMutex{}
	stringLocations     = map[string]int{}
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
	log.Println("getClientProxy", clientProxes, serverID, "=>", ret)
	clientProxiesLock.RUnlock()
	return
}

func registerClientProxyInfo(cp *ClientProxy, serverID int) {
	// register the new vacuum server client proxy
	cp.ServerID = serverID

	clientProxiesLock.Lock()
	clientProxes[serverID] = cp
	genClientProxyIDs()
	log.Printf("registerClientProxyInfo: all client proxies: %v", clientProxes)
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

	log.Printf("onClientProxyClose %v: all client proxies: %v", serverID, clientProxes)
	clientProxiesLock.Unlock()
}

func genClientProxyIDs() {
	clientProxyIDs = make([]int, 0, len(clientProxes))
	for id := range clientProxes {
		clientProxyIDs = append(clientProxyIDs, id)
	}
}

func setStringLocation(stringID string, serverID int) {
	stringLocationsLock.Lock()
	stringLocations[stringID] = serverID
	log.Printf("setStringLocation %s to %v", stringID, stringLocations)
	stringLocationsLock.Unlock()
}

func getStringLocation(stringID string) int {
	stringLocationsLock.RLock()
	log.Printf("getStringLocation %s from %v", stringID, stringLocations)
	serverID := stringLocations[stringID]
	stringLocationsLock.RUnlock()
	return serverID
}
