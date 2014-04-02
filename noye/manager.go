package noye

// Manager represents a set of managed extensions
type Manager interface {
	Respond(Message)
	Listen(IrcMessage)
	Load(string) error
	Reload(string) error
	Scripts() map[string]Script
}

// Script represents a script
type Script interface {
	Name() string
	Path() string
	Source() string
}
