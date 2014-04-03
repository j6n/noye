package irc

import "net/textproto"

// Connection is a simple connection to an irc server
type Connection struct{ conn *textproto.Conn }

// Dial uses a provided addr:port and opens a connection, returning any error
func (c *Connection) Dial(addr string) (err error) {
	conn, err := textproto.Dial("tcp", addr)
	if err != nil {
		return err
	}

	c.conn = conn
	return
}

// Close closes the connection
func (c *Connection) Close() {
	if err := c.conn.Close(); err != nil {
		log.Errorf("Closing the connection: %s\n", err)
	}
}

// WriteLine writes the 'raw' string to the connection
func (c *Connection) WriteLine(raw string) {
	if err := c.conn.Writer.PrintfLine("%s", raw); err != nil {
		log.Errorf("Writing '%s': %s\n", raw, err)
	}
}

// ReadLine returns a string, error after reading the next line from the connection
func (c *Connection) ReadLine() (string, error) {
	return c.conn.ReadLine()
}
