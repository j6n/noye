package script

import (
	"regexp"

	"github.com/robertkrimen/otto"
)

type scriptFunc func(otto.Value, ...otto.Value)

// Script represents a javascript file, with its commands/callbacks parsed.
type Script struct {
	name, path, source string

	commands  map[*regexp.Regexp]scriptFunc
	callbacks map[string][]scriptFunc
	cleanup   []scriptFunc

	subs []int64

	inits, updates map[string][]*otto.Object

	context *otto.Otto
}

func newScript(name, path, source string) *Script {
	// load extra js into this context
	return &Script{
		name: name, path: path, source: source,

		commands:  make(map[*regexp.Regexp]scriptFunc),
		callbacks: make(map[string][]scriptFunc),
		cleanup:   make([]scriptFunc, 0),

		subs: make([]int64, 0),

		inits:   make(map[string][]*otto.Object, 0),
		updates: make(map[string][]*otto.Object, 0),

		context: otto.New(),
	}
}

// Name returns the scripts name
func (s *Script) Name() string {
	return s.name
}

// Path returns the scripts path
func (s *Script) Path() string {
	return s.path
}

// Source returns the scripts source code
func (s *Script) Source() string {
	return s.source
}

// Cleanup calls all cleanup callbacks
func (s *Script) Cleanup() {
	for _, clean := range s.cleanup {
		clean(otto.NullValue())
	}
}

func (s *Script) initialize() {
	for key, vals := range s.inits {
		for _, val := range vals {
			id, ch := mq.Init(s.Name(), key, false)
			go listen(s, val, ch)
			log.Debugf("listening to '%s' %d", key, id)
			s.subs = append(s.subs, id)
		}
	}

	for key, vals := range s.updates {
		for _, val := range vals {
			id, ch := mq.Subscribe(key, false)
			go listen(s, val, ch)
			log.Debugf("listening to '%s' %d", key, id)
			s.subs = append(s.subs, id)
		}
	}
}
