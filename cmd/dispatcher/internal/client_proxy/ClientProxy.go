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
			vlog.Errorf("Client %s paniced with error: %v", cp, err)
			debug.PrintStack()
		}
	}()

	vlog.Infof("New dispatcher client: %s", cp)
	for {
		err := cp.RecvMsg(cp)
		if err != nil {
			panic(err)
		}
	}
}

func (cp *ClientProxy) HandleMsg(msg *Message, pktSize uint32, msgType MsgType_t) error {
	payload := msg[PREPAYLOAD_SIZE:pktSize]

	if msgType == CREATE_STRING_REQ {
		return cp.handleCreateStringReq(payload)
	} else if msgType == CREATE_STRING_LOCALLY_REQ {
		return cp.handleCreateStringLocallyReq(payload)
	} else if msgType == REGISTER_VACUUM_SERVER_REQ {
		return cp.handleRegisterVacuumServerReq(payload)
	} else if msgType == DECLARE_SERVICE_REQ {
		return cp.handleDeclareServiceReq(payload)
	} else if msgType == STRING_DEL_REQ {
		return cp.handleStringDelReq(payload)
	} else if msgType == START_MIGRATE_STRING_REQ {
		return cp.handleStartMigrateStringReq(payload)
	} else if msgType == MIGRATE_STRING_REQ {
		return cp.handleMigrateStringReq(payload)
	} else if msgType == LOAD_STRING_REQ {
		return cp.handleLoadStringReq(payload)
	} else {
		vlog.Panicf("ERROR: unknown dispatcher request type=%v", msgType)
		return nil
	}
}

func (cp *ClientProxy) HandleRelayMsg(msg *Message, pktSize uint32, targetStringID string) error {
	// just relay the msg
	stringInfo := getStringInfo(targetStringID)
	if !stringInfo.Migrating { // normal case
		serverID := getStringLocation(targetStringID)
		chooseServer := getClientProxy(serverID)
		vlog.Debugf(">>> RelayMsg to %s: pktSize=%v, targetID=%s", cp, chooseServer, pktSize, targetStringID)
		return chooseServer.SendAll(msg[:pktSize])
	} else { // string is migrating, we need to cache the msg until string migrated
		// just ignore for a while ...
		vlog.Debug("HandleRelayMsg: ignoring ...")
		return nil
	}
}

func (cp *ClientProxy) handleStartMigrateStringReq(data []byte) error {
	var req StartMigrateStringReq
	if err := MSG_PACKER.UnpackMsg(data, &req); err != nil {
		return err
	}

	vlog.Debugf("%s.handleStartMigrateStringReq %T %v", cp, req, req)
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

	vlog.Debugf("%s.handleMigrateStringReq %T %v", cp, req, req)

	// the string is migrating to specified server
	chooseServer := getClientProxy(req.ServerID)
	setStringLocationMigrating(req.StringID, req.ServerID, false)

	resp := MigrateStringResp{
		Name:     req.Name,
		StringID: req.StringID,
		ServerID: req.ServerID,
		Data:     req.Data,
	}

	return chooseServer.SendMsg(MIGRATE_STRING_RESP, &resp)
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

	vlog.Debugf("%s.handleCreateStringReq %T %v, choose random server: %s", cp, req, req, chooseServer)
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

	vlog.Debugf("%s.handleLoadStringReq %T %v, choose random server: %s", cp, req, req, chooseServer)
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
	vlog.Debugf("%s.handleCreateStringLocallyReq %T %v", cp, req, req)
	return nil
}

func (cp *ClientProxy) handleRegisterVacuumServerReq(data []byte) error {
	var req RegisterVacuumServerReq
	MSG_PACKER.UnpackMsg(data, &req)
	vlog.Debugf("%s.handleRegisterVacuumServerReq %T %v", cp, req, req)
	registerClientProxyInfo(cp, req.ServerID)
	return nil
}

func (cp *ClientProxy) handleDeclareServiceReq(data []byte) error {
	var req DeclareServiceReq
	if err := MSG_PACKER.UnpackMsg(data, &req); err != nil {
		return err
	}
	vlog.Debugf("%s.handleDeclareServiceReq %T %v", cp, req, req)

	// the the declare of service to all clients
	return sendToAllClientProxies(DECLARE_SERVICE_RESP, &DeclareServiceResp{
		StringID:    req.StringID,
		ServiceName: req.ServiceName,
	}, nil)
}

// String quit execution its routine on the vacuum server
func (cp *ClientProxy) handleStringDelReq(data []byte) error {
	var req StringDelReq
	if err := MSG_PACKER.UnpackMsg(data, &req); err != nil {
		return err
	}
	vlog.Debugf("%s.handleStringDelReq %T %v", cp, req, req)

	stringID := req.StringID
	return sendToAllClientProxies(STRING_DEL_RESP, &StringDelResp{
		StringID: stringID,
	}, cp) // don't send to its self
}

func sendToAllClientProxies(msgType MsgType_t, resp interface{}, exceptClientProxy *ClientProxy) error {
	clientProxiesLock.RLock()
	for _, clientProxy := range clientProxes {
		if clientProxy != exceptClientProxy {
			clientProxy.SendMsg(msgType, &resp)
		}
	}
	clientProxiesLock.RUnlock()
	return nil
}
