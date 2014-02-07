package noye

type Bot interface {
	Dial(addr, nick, user string) error
	Send(f string, a ...interface{})
	Privmsg(target, msg string)
	Join(target string)
	Part(target string)
	Close()
	Wait() <-chan struct{}
}