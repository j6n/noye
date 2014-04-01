package ext

import (
	"regexp"

	"github.com/robertkrimen/otto"
)

type scriptFunc func(otto.Value)

// Script represents a javascript file, with its commands/callbacks parsed.
type Script struct {
	Name, Path, Source string

	commands  map[*regexp.Regexp]scriptFunc
	callbacks map[string][]scriptFunc

	context *otto.Otto
}

func newScript(name, path, source string) *Script {
	return &Script{
		name, path, source,
		make(map[*regexp.Regexp]scriptFunc),
		make(map[string][]scriptFunc),
		otto.New(),
	}
}
