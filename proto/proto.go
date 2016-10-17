package proto

const (
	INVALID_MSG_TYPE           = 0
	REGISTER_VACUUM_SERVER_REQ = iota
	SEND_STRING_MESSAGE_REQ    = iota
	CREATE_STRING_REQ          = iota
	CREATE_STRING_RESP         = iota
	DECLARE_SERVICE_REQ        = iota
	DECLARE_SERVICE_RESP       = iota
)

type MsgType_t uint16

type RegisterVacuumServerReq struct {
	ServerID int `msgpack:"SID"`
}

type SendStringMessageReq struct {
	SID string
	Msg interface{}
}

type CreateStringReq struct {
	Name string `msgpack:"N"`
}

type CreateStringResp struct {
	Name string `msgpack:"N"`
}

type DeclareServiceReq struct {
	StringID    string `msgpack:"ID"`
	ServiceName string `msgpack:"SN"`
}

type DeclareServiceResp struct {
	StringID    string `msgpack:"ID"`
	ServiceName string `msgpack:"SN"`
}
