package noye

type Event interface {
	Init(Bot)
	Command() string
	Invoke(msg IrcMessage)
}
