package proto

const (
	INVALID_MSG_TYPE           = 0
	REGISTER_VACUUM_SERVER_REQ = iota
	STRING_MESSAGE_RELAY       = iota
	CREATE_STRING_REQ          = iota
	CREATE_STRING_RESP         = iota
	LOAD_STRING_REQ            = iota
	LOAD_STRING_RESP           = iota
	DECLARE_SERVICE_REQ        = iota
	DECLARE_SERVICE_RESP       = iota
	CREATE_STRING_LOCALLY_REQ  = iota
	CLOSE_STRING_RELAY         = iota
	STRING_DEL_REQ             = iota
	STRING_DEL_RESP            = iota
	MIGRATE_STRING_REQ         = iota
)

type MsgType_t uint16

type RegisterVacuumServerReq struct {
	ServerID int `msgpack:"SID"`
}

type StringMessageRelay struct {
	//StringID string      `msgpack:"ID"`
	Msg interface{} `msgpack:"M"`
}

type CreateStringReq struct {
	Name     string        `msgpack:"N"`
	StringID string        `msgpack:"ID"`
	Args     []interface{} `msgpack:"A"`
}

type CreateStringResp struct {
	Name     string        `msgpack:"N"`
	StringID string        `msgpack:"ID"`
	Args     []interface{} `msgpack:"A"`
}

type LoadStringReq struct {
	Name     string `msgpack:"N"`
	StringID string `msgpack:"ID"`
}

type LoadStringResp struct {
	Name     string `msgpack:"N"`
	StringID string `msgpack:"ID"`
}

type CreateStringLocallyReq struct {
	Name     string `msgpack:"N"`
	StringID string `msgpack:"ID"`
}

type DeclareServiceReq struct {
	StringID    string `msgpack:"ID"`
	ServiceName string `msgpack:"SN"`
}

type DeclareServiceResp struct {
	StringID    string `msgpack:"ID"`
	ServiceName string `msgpack:"SN"`
}

type StringDelReq struct {
	StringID string `msgpack:"ID"`
}

type MigrateStringReq struct {
	Name     string `msgpack:"N"`
	StringID string `msgpack:"ID"`
	ServerID int    `msgpack:"SID"`
	Data map[string]interface{}  `msgpack:"D"`
}

type StringDelResp struct {
	StringID string `msgpack:"ID"`
}

type CloseStringRelay struct {
}
