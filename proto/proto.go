package proto

const (
	INVALID_MSG_TYPE            = iota
	REGISTER_VACUUM_SERVER_REQ  = iota
	REGISTER_VACUUM_SERVER_RESP = iota
	STRING_MESSAGE_RELAY        = iota
	CREATE_STRING_REQ           = iota
	LOAD_STRING_REQ             = iota
	DECLARE_SERVICE_REQ         = iota
	CREATE_STRING_LOCALLY_REQ   = iota
	CLOSE_STRING_RELAY          = iota
	STRING_DEL_REQ              = iota
	START_MIGRATE_STRING_REQ    = iota
	MIGRATE_STRING_REQ          = iota
)

var (
	msgTypeToString = map[int]string{
		INVALID_MSG_TYPE:            "INVALID_MSG_TYPE",
		REGISTER_VACUUM_SERVER_REQ:  "REGISTER_VACUUM_SERVER_REQ",
		REGISTER_VACUUM_SERVER_RESP: "REGISTER_VACUUM_SERVER_RESP",
		STRING_MESSAGE_RELAY:        "STRING_MESSAGE_RELAY",
		CREATE_STRING_REQ:           "CREATE_STRING_REQ",
		LOAD_STRING_REQ:             "LOAD_STRING_REQ",
		DECLARE_SERVICE_REQ:         "DECLARE_SERVICE_REQ",
		CREATE_STRING_LOCALLY_REQ:   "CREATE_STRING_LOCALLY_REQ",
		CLOSE_STRING_RELAY:          "CLOSE_STRING_RELAY",
		STRING_DEL_REQ:              "STRING_DEL_REQ",
		START_MIGRATE_STRING_REQ:    "START_MIGRATE_STRING_REQ",
		MIGRATE_STRING_REQ:          "MIGRATE_STRING_REQ",
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

type LoadStringReq struct {
	Name     string        `msgpack:"N"`
	StringID string        `msgpack:"ID"`
	Args     []interface{} `msgpack:"A"`
}

type CreateStringLocallyReq struct {
	Name     string `msgpack:"N"`
	StringID string `msgpack:"ID"`
}

type DeclareServiceReq struct {
	StringID    string `msgpack:"ID"`
	ServiceName string `msgpack:"SN"`
}

type StringDelReq struct {
	StringID string `msgpack:"ID"`
}

type StartMigrateStringReq struct {
	StringID string `msgpack:"ID"`
}

type MigrateStringReq struct {
	Name             string                 `msgpack:"N"`
	StringID         string                 `msgpack:"ID"`
	ServerID         int                    `msgpack:"SID"`
	TowardsID        string                 `msgpack:"TID"`
	Args             []interface{}          `msgpack:"A"`
	Data             map[string]interface{} `msgpack:"D"`
	ExtraMigrateInfo map[string]interface{} `msgpack:"EMI"`
}

type CloseStringRelay struct {
}
