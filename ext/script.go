package ext

import (
	"regexp"

	"github.com/robertkrimen/otto"
)

type scriptFunc func(otto.Value, ...otto.Value)

// Script represents a javascript file, with its commands/callbacks parsed.
type Script struct {
	Name, Path, Source string

	commands  map[*regexp.Regexp]scriptFunc
	callbacks map[string][]scriptFunc

	context *otto.Otto
}

func newScript(name, path, source string) *Script {
	return &Script{
		Name: name, Path: path, Source: source,
		commands:  make(map[*regexp.Regexp]scriptFunc),
		callbacks: make(map[string][]scriptFunc),
		context:   otto.New(),
	}
}
