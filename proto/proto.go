package proto

const (
	INVALID_MSG_TYPE            = iota
	REGISTER_VACUUM_SERVER_REQ  = iota
	REGISTER_VACUUM_SERVER_RESP = iota
	STRING_MESSAGE_RELAY        = iota
	CREATE_STRING_REQ           = iota
	CREATE_STRING_RESP          = iota
	LOAD_STRING_REQ             = iota
	LOAD_STRING_RESP            = iota
	DECLARE_SERVICE_REQ         = iota
	DECLARE_SERVICE_RESP        = iota
	CREATE_STRING_LOCALLY_REQ   = iota
	CLOSE_STRING_RELAY          = iota
	STRING_DEL_REQ              = iota
	STRING_DEL_RESP             = iota
	START_MIGRATE_STRING_REQ    = iota
	START_MIGRATE_STRING_RESP   = iota
	MIGRATE_STRING_REQ          = iota
	MIGRATE_STRING_RESP         = iota
)

var (
	msgTypeToString = map[int]string{
		INVALID_MSG_TYPE:            "INVALID_MSG_TYPE",
		REGISTER_VACUUM_SERVER_REQ:  "REGISTER_VACUUM_SERVER_REQ",
		REGISTER_VACUUM_SERVER_RESP: "REGISTER_VACUUM_SERVER_RESP",
		STRING_MESSAGE_RELAY:        "STRING_MESSAGE_RELAY",
		CREATE_STRING_REQ:           "CREATE_STRING_REQ",
		CREATE_STRING_RESP:          "CREATE_STRING_RESP",
		LOAD_STRING_REQ:             "LOAD_STRING_REQ",
		LOAD_STRING_RESP:            "LOAD_STRING_RESP",
		DECLARE_SERVICE_REQ:         "DECLARE_SERVICE_REQ",
		DECLARE_SERVICE_RESP:        "DECLARE_SERVICE_RESP",
		CREATE_STRING_LOCALLY_REQ:   "CREATE_STRING_LOCALLY_REQ",
		CLOSE_STRING_RELAY:          "CLOSE_STRING_RELAY",
		STRING_DEL_REQ:              "STRING_DEL_REQ",
		STRING_DEL_RESP:             "STRING_DEL_RESP",
		START_MIGRATE_STRING_REQ:    "START_MIGRATE_STRING_REQ",
		START_MIGRATE_STRING_RESP:   "START_MIGRATE_STRING_RESP",
		MIGRATE_STRING_REQ:          "MIGRATE_STRING_REQ",
		MIGRATE_STRING_RESP:         "MIGRATE_STRING_RESP",
	}
)

func MsgTypeToString(msgType MsgType_t) string {
	return msgTypeToString[int(msgType)]
}

type MsgType_t uint16

type RegisterVacuumServerReq struct {
	ServerID int `msgpack:"SID"`
}

type RegisterVacuumServerResp struct {
	ServerIDS []int `msgpack:"SIDS"`
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

type StartMigrateStringReq struct {
	StringID string `msgpack:"ID"`
}

type StartMigrateStringResp struct {
	StringID string `msgpack:"ID"`
}

type MigrateStringReq struct {
	Name     string                 `msgpack:"N"`
	StringID string                 `msgpack:"ID"`
	ServerID int                    `msgpack:"SID"`
	Data     map[string]interface{} `msgpack:"D"`
}

type MigrateStringResp struct {
	Name     string                 `msgpack:"N"`
	StringID string                 `msgpack:"ID"`
	ServerID int                    `msgpack:"SID"`
	Data     map[string]interface{} `msgpack:"D"`
}

type StringDelResp struct {
	StringID string `msgpack:"ID"`
}

type CloseStringRelay struct {
}
