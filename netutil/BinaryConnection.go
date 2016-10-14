package netutil

import (
	"encoding/binary"
	"net"
)

type BinaryConnection struct {
	Connection
}

func NewBinaryConnection(conn net.Conn) BinaryConnection {
	return BinaryConnection{Connection{conn}}
}

func (bc BinaryConnection) RecvFixedLengthString(len int, pstr *string) error {
	buf := make([]byte, len)
	err := bc.RecvAll(buf)
	if err != nil {
		return err
	}
	*pstr = string(buf)
	return nil
}

func (bc BinaryConnection) SendFixedLengthString(s string) error {
	return bc.SendAll([]byte(s))
}

func (bc BinaryConnection) RecvUint16() (uint16, error) {
	buf := []byte{0, 0}
	err := bc.RecvAll(buf)
	if err != nil {
		return 0, err
	}
	return uint16(buf[0]) + (uint16(buf[1]) << 8), nil
}

func (bc BinaryConnection) SendUint16(val uint16) error {
	buf := []byte{byte(val), byte(val >> 8)}
	return bc.SendAll(buf)
}

func (bc BinaryConnection) SendInt64(val int64) error {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(val))
	return bc.SendAll(bytes)
}

//func (bc BinaryConnection) RecvSID(SID string) error {
//	err := bc.RecvFixedLengthString(SID_LENGTH, (*string)(SID))
//	return err
//}
//
//func (bc BinaryConnection) SendSID(SID string) error {
//	return bc.SendFixedLengthString(string(SID))
//}

func (bc BinaryConnection) SendString(s string) error {
	return bc.SendByteSlice([]byte(s))
}

func (bc BinaryConnection) RecvString(s *string) error {
	var buf []byte
	err := bc.RecvByteSlice(&buf)
	if err != nil {
		return err
	}
	*s = string(buf)
	return nil
}

func (bc BinaryConnection) SendByteSlice(a []byte) error {
	bc.SendUint16(uint16(len(a)))
	return bc.SendAll(a)
}

func (bc BinaryConnection) RecvByteSlice(a *[]byte) error {
	alen, err := bc.RecvUint16()
	if err != nil {
		return err
	}
	buf := make([]byte, alen)
	*a = buf
	return bc.RecvAll(buf)
}
