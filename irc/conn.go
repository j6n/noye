package irc

import "net/textproto"

type Connection struct{ conn *textproto.Conn }

func (c *Connection) Dial(addr, user string) (err error) {
	tp, err := textproto.Dial("tcp", addr)
	if err != nil {
		return err
	}

	c.conn = tp
	return
}

func (c *Connection) Close() {
	c.conn.Close()
}

func (c *Connection) WriteLine(raw string) {
	c.conn.Writer.PrintfLine("%s", raw)
}

func (c *Connection) ReadLine() (string, error) {
	return c.conn.ReadLine()
}
