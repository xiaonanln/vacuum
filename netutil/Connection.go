package netutil

import "net"

type Connection struct {
	Conn net.Conn
}

func NewConnection(conn net.Conn) Connection {
	return Connection{conn}
}

func (c Connection) String() string {
	return c.Conn.RemoteAddr().String()
}

func (c Connection) RecvByte() (byte, error) {
	buf := []byte{0}
	for {
		n, err := c.Conn.Read(buf)
		if n >= 1 {
			return buf[0], nil
		} else if err != nil {
			return 0, err
		}
	}
}

func (c Connection) SendByte(b byte) error {
	buf := []byte{b}
	for {
		n, err := c.Conn.Write(buf)
		if n >= 1 {
			return nil
		} else if err != nil {
			return err
		}
	}
}

func (c Connection) RecvAll(buf []byte) error {
	for len(buf) > 0 {
		n, err := c.Conn.Read(buf)
		if err != nil {
			return err
		}
		buf = buf[n:]
	}
	return nil
}

func (c Connection) SendAll(data []byte) error {
	for len(data) > 0 {
		n, err := c.Conn.Write(data)
		if err != nil {
			return err
		}
		data = data[n:]
	}
	return nil
}

func (c Connection) Read(data []byte) (int, error) {
	return c.Conn.Read(data)
}

func (c Connection) Write(data []byte) (int, error) {
	return c.Conn.Write(data)
}

func (c Connection) Close() error {
	return c.Conn.Close()
}
