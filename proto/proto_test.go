package proto

import (
	"log"
	"testing"
)

func TestJSONMsgPacker(t *testing.T) {
	var packer MsgPacker
	packer = &JSONMsgPacker{}
	msg := RegisterVacuumServerMsg{
		ID: "abc",
	}
	buf := make([]byte, 0, 100)
	buf, err := packer.PackMsg(msg, buf)

	log.Printf("MsgPack: %v => %s, error %v", msg, string(buf), err)
	var restoreMsg RegisterVacuumServerMsg
	err = packer.UnpackMsg(buf, &restoreMsg)
	log.Printf("MsgPack: %s => %v, error %v", string(buf), msg, err)
	if msg.ID != restoreMsg.ID {
		t.Fail()
	}
}
