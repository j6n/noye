package irc

import "net/textproto"

// Connection is a simple connection to an irc server
type Connection struct{ conn *textproto.Conn }

// Dial uses a provided addr:port and opens a connection, returning any error
func (c *Connection) Dial(addr string) (err error) {
	tp, err := textproto.Dial("tcp", addr)
	if err != nil {
		return err
	}

	c.conn = tp
	return
}

// Close closes the connection
func (c *Connection) Close() {
	c.conn.Close()
}

// WriteLine writes the 'raw' string to the connection
func (c *Connection) WriteLine(raw string) {
	c.conn.Writer.PrintfLine("%s", raw)
}

// ReadLine returns a string, error after reading the next line from the connection
func (c *Connection) ReadLine() (string, error) {
	return c.conn.ReadLine()
}
