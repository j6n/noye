package noye

// Plugin is an abstract contract for a Bot plugin
type Plugin interface {
	Hook(Bot)
	Listen() chan Message
	Name() string
	Status(string) bool
	SetStatus(string, bool)
}
