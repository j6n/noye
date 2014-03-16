package noye

// Event is an abstract contract for an IRC event
type Event interface {
	Init(Bot)
	Command() string
	Invoke(msg IrcMessage)
}
