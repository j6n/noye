package noye

type Manager interface {
	Respond(Message)
	Listen(IrcMessage)
	Load(string) error
	Reload(string) error
}
