package proto

const (
	REGISTER_VACUUM_SERVER_REQ = 1
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
