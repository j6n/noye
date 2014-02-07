package noye

type Bot interface {
	Dial(addr, nick, user string) error
	Close()

	Send(f string, a ...interface{})
	Privmsg(target, msg string)

	Join(target string)
	Part(target string)

	Wait() <-chan struct{}
	AddPlugin(Plugin)
}
