package noye

// Conn is an abstract interface for an IRC connection
type Conn interface {
	Dial(addr string) error
	Close()

	WriteLine(raw string)
	ReadLine() (string, error)
}
