package noye

type Conn interface {
	Dial(addr, user string) error
	Close()

	WriteLine(raw string)
	ReadLine() (string, error)
}
