package dispatcher_client

import (
	"log"

	"time"

	"sync/atomic"

	"unsafe"

	"errors"

	"runtime/debug"

	"github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/msgbufpool"
	"github.com/xiaonanln/vacuum/netutil"
	"github.com/xiaonanln/vacuum/proto"
)

const (
	LOOP_DELAY_ON_DISPATCHER_CLIENT_ERROR = 3 * time.Second
)

var (
	_dispatcherClient         *DispatcherClient
	serverID                  = 0
	errDispatcherNotConnected = errors.New("dispatcher not connected")
)

func getDispatcherClient() *DispatcherClient {
	addr := (*uintptr)(unsafe.Pointer(&_dispatcherClient))
	return (*DispatcherClient)(unsafe.Pointer(atomic.LoadUintptr(addr)))
}

func setDispatcherClient(dc *DispatcherClient) {
	addr := (*uintptr)(unsafe.Pointer(&_dispatcherClient))
	atomic.StoreUintptr(addr, uintptr(unsafe.Pointer(dc)))
}

func assureConnectedDispatcherClient() *DispatcherClient {
	var err error
	dispatcherClient := getDispatcherClient()
	log.Println("assureConnectedDispatcherClient: dispatcherClient", dispatcherClient)
	for dispatcherClient == nil {
		dispatcherClient, err = connectDispatchClient()
		if err != nil {
			log.Printf("Connect to dispatcher failed: %s", err.Error())
			time.Sleep(LOOP_DELAY_ON_DISPATCHER_CLIENT_ERROR)
			continue
		}

		if serverID == 0 {
			log.Panicf("invalid serverID: %v", serverID)
		}

		dispatcherClient.RegisterVacuumServer(serverID)
		setDispatcherClient(dispatcherClient)
	}

	return dispatcherClient
}

func connectDispatchClient() (*DispatcherClient, error) {
	conn, err := netutil.ConnectTCP("localhost", 7581)
	if err != nil {
		return nil, err
	}
	return newDispatcherClient(conn), nil
}

func Initialize(_serverID int, h DispatcherRespHandler) {
	serverID = _serverID
	dispatcherRespHandler = h

	assureConnectedDispatcherClient()
	go netutil.ServeForever(serveDispatcherClient)
}

func SendStringMessage(stringID string, msg common.StringMessage) error {
	dispatcherClient := getDispatcherClient()
	if dispatcherClient == nil {
		debug.PrintStack()
		log.Printf("dispatcher client is nil")
		return errDispatcherNotConnected
	}
	return dispatcherClient.SendStringMessage(stringID, msg)
}

func SendCreateStringReq(name string, stringID string) error {
	dispatcherClient := getDispatcherClient()
	if dispatcherClient == nil {
		debug.PrintStack()
		log.Printf("dispatcher client is nil")
		return errDispatcherNotConnected
	}
	return dispatcherClient.SendCreateStringReq(name, stringID)
}

func SendCreateStringLocallyReq(name string, stringID string) error {
	dispatcherClient := getDispatcherClient()
	if dispatcherClient == nil {
		debug.PrintStack()
		log.Printf("dispatcher client is nil")
		return errDispatcherNotConnected
	}
	return dispatcherClient.SendCreateStringLocallyReq(name, stringID)
}

func SendDeclareServiceReq(sid string, serviceName string) error {
	dispatcherClient := getDispatcherClient()
	if dispatcherClient == nil {
		debug.PrintStack()
		log.Printf("dispatcher client is nil")
		return errDispatcherNotConnected
	}
	return dispatcherClient.SendDeclareServiceReq(sid, serviceName)
}

// serve the dispatcher client, receive RESPs from dispatcher and process
func serveDispatcherClient() {
	var err error
	log.Printf("serveDispatcherClient: start serving dispatcher client ...")
	for {
		dispatcherClient := assureConnectedDispatcherClient()

		var msgPackInfo proto.MsgPacketInfo
		err = dispatcherClient.RecvMsgPacket(&msgPackInfo)
		if err != nil {
			log.Printf("serveDispatcherClient: RecvMsgPacket error: %s", err.Error())
			dispatcherClient.Close()
			setDispatcherClient(nil)
			time.Sleep(LOOP_DELAY_ON_DISPATCHER_CLIENT_ERROR)
			continue
		}

		log.Printf("serveDispatcherClient: received dispatcher resp: %v", msgPackInfo)
		// handle the packet ... on this vacuum server
		msgtype := msgPackInfo.MsgType
		if msgtype == proto.SEND_STRING_MESSAGE_RESP {
			// receive string message.
			err = handleSendStringMessageResp(dispatcherClient, msgPackInfo.Payload)
		} else if msgtype == proto.CREATE_STRING_RESP {
			// create real string instance
			err = handleCreateStringResp(dispatcherClient, msgPackInfo.Payload)
		} else if msgtype == proto.DECLARE_SERVICE_RESP {
			// declare service
			err = handleDeclareServiceResp(dispatcherClient, msgPackInfo.Payload)
		} else {
			log.Panicf("serveDispatcherClient: invalid msg type: %v", msgtype)
		}

		// reclaim the msgbuf
		msgbufpool.PutMsgBuf(msgPackInfo.Msgbuf)
	}
}

func handleSendStringMessageResp(_ *DispatcherClient, payload []byte) error {
	var resp proto.SendStringMessageResp
	err := proto.MSG_PACKER.UnpackMsg(payload, &resp)
	if err != nil {
		return err
	}

	dispatcherRespHandler.HandleDispatcherResp_SendStringMessage(resp.StringID, resp.Msg)
	return nil
}

func handleCreateStringResp(_ *DispatcherClient, payload []byte) error {
	var resp proto.CreateStringResp
	err := proto.MSG_PACKER.UnpackMsg(payload, &resp)
	if err != nil {
		return err
	}

	dispatcherRespHandler.HandleDispatcherResp_CreateString(resp.Name, resp.StringID)
	return nil
}

func handleDeclareServiceResp(_ *DispatcherClient, payload []byte) error {
	var resp proto.DeclareServiceResp
	err := proto.MSG_PACKER.UnpackMsg(payload, &resp)
	if err != nil {
		return err
	}

	dispatcherRespHandler.HandleDispatcherResp_DeclareService(resp.StringID, resp.ServiceName)
	return nil
}
