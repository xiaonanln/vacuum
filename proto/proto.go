package proto

const (
	REGISTER_VACUUM_SERVER_REQ = 1
)

type msgtype_t uint16

type RegisterVacuumServerReq struct {
	ID string
}
