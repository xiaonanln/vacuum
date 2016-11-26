package client_proxy

import (
	"net"

	"fmt"

	"runtime/debug"

	"github.com/xiaonanln/vacuum/netutil"
	. "github.com/xiaonanln/vacuum/proto"
	"github.com/xiaonanln/vacuum/vlog"
)

type ClientProxy struct {
	MessageConnection
	ServerID int
}

func NewClientProxy(conn net.Conn) *ClientProxy {
	return &ClientProxy{
		MessageConnection: NewMessageConnection(conn),
		ServerID:          0, // to be registered
	}
}

func (cp *ClientProxy) String() string {
	return fmt.Sprintf("Server<%d>", cp.ServerID)
}

func (cp *ClientProxy) Serve() {
	defer func() {
		cp.Close()
		onClientProxyClose(cp)

		err := recover()

		if err != nil && !netutil.IsConnectionClosed(err) {
			vlog.Error("Client %s paniced with error: %v", cp, err)
			debug.PrintStack()
		}
	}()

	vlog.Info("New dispatcher client: %s", cp)
	for {
		err := cp.RecvMsg(cp)
		if err != nil {
			panic(err)
		}
	}
}

func (cp *ClientProxy) HandleMsg(msg *Message, pktSize uint32, msgType MsgType_t) error {
	payload := msg[PREPAYLOAD_SIZE:pktSize]

	var err error
	if msgType == CREATE_STRING_REQ {
		err = cp.handleCreateStringReq(payload)
	} else if msgType == CREATE_STRING_LOCALLY_REQ {
		err = cp.handleCreateStringLocallyReq(payload)
	} else if msgType == REGISTER_VACUUM_SERVER_REQ {
		err = cp.handleRegisterVacuumServerReq(payload)
	} else if msgType == DECLARE_SERVICE_REQ {
		err = cp.handleDeclareServiceReq(payload)
	} else if msgType == STRING_DEL_REQ {
		err = cp.handleStringDelReq(payload)
	} else if msgType == START_MIGRATE_STRING_REQ {
		err = cp.handleStartMigrateStringReq(payload)
	} else if msgType == MIGRATE_STRING_REQ {
		err = cp.handleMigrateStringReq(payload)
	} else if msgType == LOAD_STRING_REQ {
		err = cp.handleLoadStringReq(payload)
	} else {
		vlog.Panicf("ERROR: unknown dispatcher request type=%v", msgType)
	}

	msg.Release()

	return err
}

func (cp *ClientProxy) HandleRelayMsg(msg *Message, pktSize uint32, targetStringID string) (err error) {
	// just relay the msg
	stringCtrl := getStringCtrl(targetStringID) // TODO: optimize: lock for read first
	stringCtrl.RLock()

	if !stringCtrl.Migrating { // normal case
		serverID := stringCtrl.ServerID

		chooseServer := getClientProxy(serverID)
		vlog.Debug(">>> RelayMsg to %s: pktSize=%v, targetID=%s", chooseServer, pktSize, targetStringID)
		err = chooseServer.SendAll(msg[:pktSize])
		// FIXME: the write lock here affect the order or relay messages, we need to use better lock control
		stringCtrl.RUnlock()

		msg.Release()
		return
	}

	// string is migrating ...
	stringCtrl.RUnlock()
	stringCtrl.Lock() // re-lock for write
	defer stringCtrl.Unlock()

	if !stringCtrl.Migrating { // normal case
		serverID := stringCtrl.ServerID

		chooseServer := getClientProxy(serverID)
		vlog.Debug(">>> RelayMsg to %s: pktSize=%v, targetID=%s", chooseServer, pktSize, targetStringID)
		err = chooseServer.SendAll(msg[:pktSize])
		msg.Release()
		return
	}

	stringCtrl.cachedMessages = append(stringCtrl.cachedMessages, _CachedMessage{
		msg:     msg,
		pktsize: pktSize,
	}) // cahce the message
	// not release the msg

	return
}

func (cp *ClientProxy) handleStartMigrateStringReq(data []byte) error {
	var req StartMigrateStringReq
	if err := MSG_PACKER.UnpackMsg(data, &req); err != nil {
		return err
	}

	vlog.Debug("%s.handleStartMigrateStringReq %T %v", cp, req, req)
	// migrating, messages to this String should be cached, until real migration happened
	setStringMigrating(req.StringID, true)
	// send the resp to the client
	resp := StartMigrateStringResp{
		StringID: req.StringID,
	}

	return cp.SendMsg(START_MIGRATE_STRING_RESP, resp)
}

func (cp *ClientProxy) handleMigrateStringReq(data []byte) error {
	var req MigrateStringReq
	if err := MSG_PACKER.UnpackMsg(data, &req); err != nil {
		return err
	}

	vlog.Debug("%s.handleMigrateStringReq %T %v", cp, req, req)

	// the string is migrating to specified server
	chooseServer := getClientProxy(req.ServerID)

	ctrl := getStringCtrl(req.StringID)
	ctrl.Lock()
	defer ctrl.Unlock()

	ctrl.ServerID = req.ServerID
	ctrl.Migrating = false

	var cacheMessages []_CachedMessage
	cacheMessages, ctrl.cachedMessages = ctrl.cachedMessages, nil

	resp := MigrateStringResp{
		Name:     req.Name,
		StringID: req.StringID,
		ServerID: req.ServerID,
		Args:     req.Args,
		Data:     req.Data,
	}

	if err := chooseServer.SendMsg(MIGRATE_STRING_RESP, &resp); err != nil {
		return err
	}

	for _, cachedMsg := range cacheMessages {
		if err := chooseServer.SendAll(cachedMsg.msg[:cachedMsg.pktsize]); err != nil {
			return err
		}
		cachedMsg.msg.Release()
	}

	return nil
}

func (cp *ClientProxy) handleCreateStringReq(data []byte) error {
	var req CreateStringReq
	if err := MSG_PACKER.UnpackMsg(data, &req); err != nil {
		return err
	}

	// choose one server for create string

	chooseServer := getRandomClientProxy()
	// save the stringID with the serverID
	stringID := req.StringID
	setStringLocation(stringID, chooseServer.ServerID)

	vlog.Debug("%s.handleCreateStringReq %T %v, choose random server: %s", cp, req, req, chooseServer)
	resp := CreateStringResp{
		Name:     req.Name,
		StringID: stringID,
		Args:     req.Args,
	}

	return chooseServer.SendMsg(CREATE_STRING_RESP, &resp)
}

func (cp *ClientProxy) handleLoadStringReq(data []byte) error {
	var req LoadStringReq
	if err := MSG_PACKER.UnpackMsg(data, &req); err != nil {
		return err
	}

	chooseServer := getRandomClientProxy()
	stringID := req.StringID
	setStringLocation(stringID, chooseServer.ServerID)

	vlog.Debug("%s.handleLoadStringReq %T %v, choose random server: %s", cp, req, req, chooseServer)
	resp := LoadStringResp{
		Name:     req.Name,
		StringID: stringID,
	}

	return chooseServer.SendMsg(LOAD_STRING_RESP, &resp)
}

func (cp *ClientProxy) handleCreateStringLocallyReq(data []byte) error {
	var req CreateStringLocallyReq
	if err := MSG_PACKER.UnpackMsg(data, &req); err != nil {
		return err
	}

	// choose one server for create string

	stringID := req.StringID
	setStringLocation(stringID, cp.ServerID)
	vlog.Debug("%s.handleCreateStringLocallyReq %T %v", cp, req, req)
	return nil
}

func (cp *ClientProxy) handleRegisterVacuumServerReq(data []byte) error {
	var req RegisterVacuumServerReq
	MSG_PACKER.UnpackMsg(data, &req)
	vlog.Debug("%s.handleRegisterVacuumServerReq %T %v", cp, req, req)
	serverIDs := registerClientProxyInfo(cp, req.ServerID) // register the game server

	return sendToAllClientProxies(REGISTER_VACUUM_SERVER_RESP, &RegisterVacuumServerResp{ServerIDS: serverIDs})
}

func (cp *ClientProxy) handleDeclareServiceReq(data []byte) error {
	var req DeclareServiceReq
	if err := MSG_PACKER.UnpackMsg(data, &req); err != nil {
		return err
	}
	vlog.Debug("%s.handleDeclareServiceReq %T %v", cp, req, req)

	// the the declare of service to all clients
	return sendToAllClientProxies(DECLARE_SERVICE_RESP, &DeclareServiceResp{
		StringID:    req.StringID,
		ServiceName: req.ServiceName,
	})
}

// String quit execution its routine on the vacuum server
func (cp *ClientProxy) handleStringDelReq(data []byte) error {
	var req StringDelReq
	if err := MSG_PACKER.UnpackMsg(data, &req); err != nil {
		return err
	}
	vlog.Debug("%s.handleStringDelReq %T %v", cp, req, req)

	stringID := req.StringID
	return sendToAllClientProxies(STRING_DEL_RESP, &StringDelResp{
		StringID: stringID,
	}) // don't send to its self
}

func sendToAllClientProxies(msgType MsgType_t, resp interface{}) error {
	return sendToAllClientProxiesExcept(msgType, resp, nil)
}

func sendToAllClientProxiesExcept(msgType MsgType_t, resp interface{}, exceptClientProxy *ClientProxy) error {
	clientProxiesLock.RLock()
	for _, clientProxy := range clientProxes {
		if clientProxy != exceptClientProxy {
			clientProxy.SendMsg(msgType, &resp)
		}
	}
	clientProxiesLock.RUnlock()
	return nil
}
