package dispatcher_client

import (
	"log"

	"time"

	"sync/atomic"

	"unsafe"

	"github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/msgbufpool"
	"github.com/xiaonanln/vacuum/netutil"
	"github.com/xiaonanln/vacuum/proto"
)

const (
	LOOP_DELAY_ON_DISPATCHER_CLIENT_ERROR = 3 * time.Second
)

var (
	_dispatcherClient *DispatcherClient
	serverID          = 0
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
	log.Println("dispatcherClient", dispatcherClient)
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

func Initialize(_serverID int) {
	serverID = _serverID
	go netutil.ServeForever(serveDispatcherClient)
	assureConnectedDispatcherClient()
}

func SendStringMessage(sid string, msg common.StringMessage) {
	var err error
	dispatcherClient := assureConnectedDispatcherClient()
	err = dispatcherClient.SendStringMessage(sid, msg)
	if err != nil {
		log.Printf("SendStringMessage: send string message failed with error %s, dispatcher lost ..", err.Error())
		dispatcherClient.Close()
		setDispatcherClient(nil)
	}
}

func CreateString(name string) error {
	dispatcherClient := assureConnectedDispatcherClient()
	return dispatcherClient.CreateString(name)
}

// serve the dispatcher client, receive RESPs from dispatcher and process
func serveDispatcherClient() {
	var err error
	log.Printf("serveDispatcherClient: start serving dispatcher client ...")
	for {
		dispatcherClient := getDispatcherClient()
		if dispatcherClient == nil {
			log.Printf("serveDispatcherClient: dispatcher client is nil")
			time.Sleep(LOOP_DELAY_ON_DISPATCHER_CLIENT_ERROR)
			continue
		}
		var msgPackInfo proto.MsgPacketInfo
		err = dispatcherClient.RecvMsgPacket(&msgPackInfo)
		if err != nil {
			log.Printf("serveDispatcherClient: RecvMsgPacket error: %s", err.Error())
			time.Sleep(LOOP_DELAY_ON_DISPATCHER_CLIENT_ERROR)
			continue
		}

		log.Printf("serveDispatcherClient: received dispatcher resp: %v", msgPackInfo)

		// handle the packet ...
		msgtype := msgPackInfo.MsgType
		if msgtype == proto.CREATE_STRING_RESP {
			// create string on this vacuum server
			err = handleCreateStringResp(dispatcherClient, msgPackInfo.Payload)
		} else {
			log.Panicf("serveDispatcherClient: invalid msg type: %v", msgtype)
		}

		// reclaim the msgbuf
		msgbufpool.PutMsgBuf(msgPackInfo.Msgbuf)
	}
}

func handleCreateStringResp(dispatcherClient *DispatcherClient, payload []byte) error {
	var resp proto.CreateStringResp
	err := proto.MSG_PACKER.UnpackMsg(payload, &resp)
	if err != nil {
		return err
	}

	//vacuum.createString(resp.Name)
	return nil
}
