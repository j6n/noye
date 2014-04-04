package ext

import (
	"fmt"
	"regexp"

	"github.com/robertkrimen/otto"
)

type scriptFunc func(otto.Value, ...otto.Value)

// Script represents a javascript file, with its commands/callbacks parsed.
type Script struct {
	name, path, source string

	commands  map[*regexp.Regexp]scriptFunc
	callbacks map[string][]scriptFunc

	subs  []int64
	inits []otto.Value

	context *otto.Otto
}

func newScript(name, path, source string) *Script {
	context := otto.New()
	if data, err := lodashminjs(); err == nil {
		context.Run(string(data))
	}

	return &Script{
		name: name, path: path, source: source,
		commands:  make(map[*regexp.Regexp]scriptFunc),
		callbacks: make(map[string][]scriptFunc),
		inits:     make([]otto.Value, 0),
		subs:      make([]int64, 0),
		context:   context,
	}
}

func (s *Script) String() string {
	return fmt.Sprintf("%s @ %s", s.Name, s.Path)
}

// Name returns the scripts name
func (s *Script) Name() string { return s.name }

// Path returns the scripts path
func (s *Script) Path() string { return s.path }

// Source returns the scripts source code
func (s *Script) Source() string { return s.source }
