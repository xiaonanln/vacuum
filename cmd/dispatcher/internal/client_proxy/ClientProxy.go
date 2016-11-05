package client_proxy

import (
	"net"

	log "github.com/Sirupsen/logrus"

	"fmt"

	"runtime/debug"

	. "github.com/xiaonanln/vacuum/proto"
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
		if err != nil {
			log.Errorf("Client %s paniced with error: %v", cp, err)
			debug.PrintStack()
		}
	}()

	log.Infof("New dispatcher client: %s", cp)
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
		cp.handleCreateStringReq(payload)
	} else if msgType == CREATE_STRING_LOCALLY_REQ {
		cp.handleCreateStringLocallyReq(payload)
	} else if msgType == REGISTER_VACUUM_SERVER_REQ {
		cp.handleRegisterVacuumServerReq(payload)
	} else if msgType == DECLARE_SERVICE_REQ {
		cp.handleDeclareServiceReq(payload)
	} else if msgType == STRING_DEL_REQ {
		cp.handleStringDelReq(payload)
	} else if msgType == MIGRATE_STRING_REQ {
		cp.handleMigrateStringReq(payload)
	} else if msgType == LOAD_STRING_REQ {
		cp.handleLoadStringReq(payload)
	} else {
		log.Panicf("ERROR: unknown dispatcher request type=%v", msgType)
	}
	return nil
}

func (cp *ClientProxy) HandleRelayMsg(msg *Message, pktSize uint32, targetID string) error {
	// just relay the msg
	serverID := getStringLocation(targetID)
	chooseServer := getClientProxy(serverID)
	log.WithFields(log.Fields{"pktSize": pktSize, "targetID": targetID}).Debugf("%s.HandleRelayMsg to %s", cp, chooseServer)
	return chooseServer.SendAll(msg[:pktSize])
}

//func (cp *ClientProxy) handleSendStringMessageReq(data []byte) {
//	var req SendStringMessageReq
//	MSG_PACKER.UnpackMsg(data, &req)
//
//	targetStringID := req.StringID
//	resp := SendStringMessageResp{
//		StringID: targetStringID,
//		Msg:      req.Msg,
//	}
//
//	serverID := getStringLocation(targetStringID)
//	chooseServer := getClientProxy(serverID)
//
//	log.Debugf("%s.handleSendStringMessageReq %T %v, target server %s", cp, req, req, chooseServer)
//	chooseServer.SendMsg(SEND_STRING_MESSAGE_RESP, &resp)
//}

func (cp *ClientProxy) handleCreateStringReq(data []byte) {
	var req CreateStringReq
	MSG_PACKER.UnpackMsg(data, &req)

	// choose one server for create string

	chooseServer := getRandomClientProxy()
	// save the stringID with the serverID
	stringID := req.StringID
	setStringLocation(stringID, chooseServer.ServerID)

	log.Debugf("%s.handleCreateStringReq %T %v, choose random server: %s", cp, req, req, chooseServer)
	resp := CreateStringResp{
		Name:     req.Name,
		StringID: stringID,
		Args:     req.Args,
	}

	chooseServer.SendMsg(CREATE_STRING_RESP, &resp)
}

func (cp *ClientProxy) handleLoadStringReq(data []byte) {
	var req LoadStringReq
	MSG_PACKER.UnpackMsg(data, &req)

	chooseServer := getRandomClientProxy()
	stringID := req.StringID
	setStringLocation(stringID, chooseServer.ServerID)

	log.Debugf("%s.handleLoadStringReq %T %v, choose random server: %s", cp, req, req, chooseServer)
	resp := LoadStringResp{
		Name:     req.Name,
		StringID: stringID,
	}

	chooseServer.SendMsg(LOAD_STRING_RESP, &resp)
}

func (cp *ClientProxy) handleCreateStringLocallyReq(data []byte) {
	var req CreateStringLocallyReq
	MSG_PACKER.UnpackMsg(data, &req)

	// choose one server for create string

	stringID := req.StringID
	setStringLocation(stringID, cp.ServerID)
	log.Debugf("%s.handleCreateStringLocallyReq %T %v", cp, req, req)
}

func (cp *ClientProxy) handleRegisterVacuumServerReq(data []byte) {
	var req RegisterVacuumServerReq
	MSG_PACKER.UnpackMsg(data, &req)
	log.Debugf("%s.handleRegisterVacuumServerReq %T %v", cp, req, req)
	registerClientProxyInfo(cp, req.ServerID)
}

func (cp *ClientProxy) handleDeclareServiceReq(data []byte) {
	var req DeclareServiceReq
	MSG_PACKER.UnpackMsg(data, &req)
	log.Debugf("%s.handleDeclareServiceReq %T %v", cp, req, req)

	// the the declare of service to all clients
	sendToAllClientProxies(DECLARE_SERVICE_RESP, &DeclareServiceResp{
		StringID:    req.StringID,
		ServiceName: req.ServiceName,
	}, nil)
}

// String quit execution its routine on the vacuum server
func (cp *ClientProxy) handleStringDelReq(data []byte) {
	var req StringDelReq
	MSG_PACKER.UnpackMsg(data, &req)
	log.Debugf("%s.handleStringDelReq %T %v", cp, req, req)

	stringID := req.StringID
	sendToAllClientProxies(STRING_DEL_RESP, &StringDelResp{
		StringID: stringID,
	}, cp) // don't send to its self
}

func (cp *ClientProxy) handleMigrateStringReq(data []byte) {
	var req MigrateStringReq
	MSG_PACKER.UnpackMsg(data, &req)
	log.Debugf("%s.handleMigrateStringReq %T %v", cp, req, req)

	// the string is migrating to specified server
}

func sendToAllClientProxies(msgType MsgType_t, resp interface{}, exceptClientProxy *ClientProxy) {
	clientProxiesLock.RLock()
	for _, clientProxy := range clientProxes {
		if clientProxy != exceptClientProxy {
			clientProxy.SendMsg(msgType, &resp)
		}
	}
	clientProxiesLock.RUnlock()
}
