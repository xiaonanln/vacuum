package proto

const (
	INVALID_MSG_TYPE           = 0
	REGISTER_VACUUM_SERVER_REQ = iota
	SEND_STRING_MESSAGE_REQ    = iota
)

type MsgType_t uint16

type RegisterVacuumServerReq struct {
	ID string
}

type SendStringMessageReq struct {
	SID string
	Msg interface{}
}
