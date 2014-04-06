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
