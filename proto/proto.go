package proto

const (
	REGISTER_VACUUM_SERVER_REQ = 1
)

type MsgType_t uint16

type RegisterVacuumServerReq struct {
	ID string
}
