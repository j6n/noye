package noye

// Manager represents a set of managed extensions
type Manager interface {
	Respond(Message)
	Listen(IrcMessage)

	LoadScripts(dir string)
	Load(string) error

	Reload(string) error
	ReloadBase() error

	Scripts() []Script
}

// Script represents a script
type Script interface {
	Name() string
	Path() string
	Source() string
}
