package noye

// Script represents a script
type Script interface {
	Name() string
	Path() string
	Source() string
	Cleanup()
}
