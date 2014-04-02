package noye

// Manager represents a set of managed extensions
type Manager interface {
	Respond(Message)
	Listen(IrcMessage)
	Load(string) error
	Reload(string) error
	Eval(source string) error
}
