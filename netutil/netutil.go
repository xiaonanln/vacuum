package netutil

import (
	"fmt"
	"net"
	"runtime"
)

var (
	NEWLINE = "\n"
)

func init() {
	if runtime.GOOS == "windows" {
		NEWLINE = "\n\r"
	}
}

func IsTemporaryNetError(err error) bool {
	if err == nil {
		return false
	}

	netErr, ok := err.(net.Error)
	if !ok {
		return false
	}
	return netErr.Temporary() || netErr.Timeout()
}

func WriteAll(conn net.Conn, data []byte) error {
	for len(data) > 0 {
		n, err := conn.Write(data)
		if n > 0 {
			data = data[n:]
		}
		if err != nil {
			if IsTemporaryNetError(err) {
				continue
			} else {
				return err
			}
		}
	}
	return nil
}

func ReadAll(conn net.Conn, data []byte) error {
	for len(data) > 0 {
		n, err := conn.Read(data)
		if n > 0 {
			data = data[n:]
		}
		if err != nil {
			if IsTemporaryNetError(err) {
				continue
			} else {
				return err
			}
		}
	}
	return nil
}

func ReadLine(conn net.Conn) (string, error) {
	var _linebuff [1024]byte
	linebuff := _linebuff[0:0]

	buff := [1]byte{0} // buff of just 1 byte

	for {
		n, err := conn.Read(buff[0:1])
		if err != nil {
			if IsTemporaryNetError(err) {
				continue
			} else {
				return "", err
			}
		}
		if n == 1 {
			c := buff[0]
			if c == '\n' {
				return string(linebuff), nil
			} else {
				linebuff = append(linebuff, c)
			}
		}
	}
}

func ConnectTCP(host string, port int) (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr)
	return conn, err
}
