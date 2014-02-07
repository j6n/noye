package irc

import "net/textproto"

type Conn interface {
	Dial(addr, user string) error
	Close()
	WriteLine(raw string)
	ReadLine() (string, error)
}

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

/*
	c.WriteLine(fmt.Sprintf("NICK %s", user))
	c.WriteLine(fmt.Sprintf("USER %s * 0 :%s", user, user))
*/
