package noye

// Manager represents a set of managed scripts
type Manager interface {
	Respond(msg Message)
	Listen(msg IrcMessage)

	LoadScripts(dir string)
	LoadFile(path string) error

	Reload(script string) error
	ReloadBase() error

	Unload(name string) error
	UnloadAll()

	Scripts() map[string]Script
}
