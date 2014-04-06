package noye

// Bot is an abstract contract for an IRC Bot
type Bot interface {
	Dial(addr, nick, user string) error
	Close()

	Send(f string, a ...interface{})
	Privmsg(target, msg string)

	Join(target string)
	Part(target string)

	Quit()

	Wait() <-chan struct{}
	Ready() <-chan struct{}

	Manager() Manager
}
