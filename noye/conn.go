package noye

type Conn interface {
	Dial(addr string) error
	Close()

	WriteLine(raw string)
	ReadLine() (string, error)
}
