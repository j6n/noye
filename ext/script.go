package ext

import (
	"fmt"
	"regexp"

	"github.com/j6n/noye/etc/html"

	"github.com/j6n/noye/etc/http"
	"github.com/j6n/noye/store"
	"github.com/robertkrimen/otto"
)

type scriptFunc func(otto.Value, ...otto.Value)

// Script represents a javascript file, with its commands/callbacks parsed.
type Script struct {
	name, path, source string

	commands  map[*regexp.Regexp]scriptFunc
	callbacks map[string][]scriptFunc
	cleanup   []scriptFunc

	subs    []int64
	updates []inits

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
		cleanup:   make([]scriptFunc, 0),

		subs:    make([]int64, 0),
		updates: make([]inits, 0),

		context: context,
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

// Cleanup calls all cleanup callbacks
func (s *Script) Cleanup() {
	for _, clean := range s.cleanup {
		clean(otto.NullValue())
	}
}

// these are the default methods injected into it
// set saves a string for a key
func (s *Script) scriptSet(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) != 2 || !call.ArgumentList[0].IsString() || !call.ArgumentList[1].IsString() {
		return otto.FalseValue()
	}

	key, data := call.ArgumentList[0].String(), call.ArgumentList[1].String()
	if err := store.Set(s.Name(), key, data); err != nil {
		log.Errorf("(%s) setting val at '%s': %s", s.Name(), key, err)
		return otto.FalseValue()
	}

	return otto.TrueValue()
}

// set gets a string for a key
func (s *Script) scriptGet(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) != 1 || !call.ArgumentList[0].IsString() {
		return otto.UndefinedValue()
	}

	key := call.ArgumentList[0].String()
	data, err := store.Get(s.Name(), key)
	if err != nil {
		log.Errorf("(%s) getting val at '%s': %s", s.Name(), key, err)
		return otto.UndefinedValue()
	}

	val, _ := s.context.ToValue(data)
	return val
}

type inits struct {
	table string
	init  bool
	obj   *otto.Object
}

// subInit subs and fills channel with default values
func (s *Script) scriptSubInit(init bool) func(otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 2 || !call.ArgumentList[0].IsString() || !call.ArgumentList[1].IsFunction() {
			return otto.FalseValue()
		}

		key, fn := call.ArgumentList[0].String(), call.ArgumentList[1].Object()
		s.updates = append(s.updates, inits{key, init, fn})
		return otto.TrueValue()
	}
}

func (s *Script) doInits() {
	for _, v := range s.updates {
		fn := v.obj
		var ch chan string
		var id int64

		if v.init {
			id, ch = mq.Init(s.Name(), v.table, false)
		} else {
			id, ch = mq.Subscribe(v.table, false)
		}

		go func(o *otto.Object) {
			fn := o.Value()
			for data := range ch {
				val, err := s.context.ToValue(data)
				if err != nil {
					log.Errorf("(%s) convert val '%s': %s", s.Name(), data, err)
					continue
				}

				if _, err = fn.Call(fn, val); err != nil {
					log.Errorf("(%s) calling fn: %s", s.Name(), err)
				}
			}
		}(fn)

		s.subs = append(s.subs, id)
	}
}

// unsub unsubscribes to the message queue
func (s *Script) scriptUnsub(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) < 1 || !call.ArgumentList[0].IsNumber() {
		return otto.FalseValue()
	}
	id, err := call.ArgumentList[0].ToInteger()
	if err != nil {
		log.Errorf("(%s) wasn't given an unsub id: %s", s.Name(), err)
		return otto.TrueValue()
	}

	mq.Unsubscribe(id)
	for i, sub := range s.subs {
		if sub == id {
			s.subs = s.subs[:i+copy(s.subs[i:], s.subs[i+1:])]
		}
	}
	return otto.TrueValue()
}

// update broadcasts a message for a key
func (s *Script) scriptUpdate(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) != 2 || !call.ArgumentList[0].IsString() || !call.ArgumentList[1].IsString() {
		return otto.FalseValue()
	}

	key, val := call.ArgumentList[0].String(), call.ArgumentList[1].String()
	mq.Update(key, val, false)
	return otto.TrueValue()
}

// httpget does a httpget for a string, returning a int, string
func (s *Script) scriptHttpget(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) < 2 || !call.ArgumentList[0].IsString() || !call.ArgumentList[1].IsFunction() {
		return otto.FalseValue()
	}

	headers := make(map[string]string)
	if len(call.ArgumentList) > 2 {
		obj, err := call.ArgumentList[2].Export()
		if err != nil {
			otto.FalseValue()
		}

		if m, ok := obj.(map[string]interface{}); ok {
			for k, v := range m {
				headers[k] = v.(string)
			}
		}
	}

	url, fn := call.ArgumentList[0].String(), call.ArgumentList[1]
	go func() {
		res, status := http.Get(url, headers)
		sval, _ := s.context.ToValue(status)
		rval, _ := s.context.ToValue(res)
		fn.Call(fn, sval, rval)
	}()
	return otto.TrueValue()
}

func (s *Script) scriptHttpfollow(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) < 2 || !call.Argument(0).IsString() || !call.Argument(1).IsFunction() {
		return otto.FalseValue()
	}

	url, fn := call.Argument(0).String(), call.Argument(1)
	go func() {
		res, status := http.Follow(url)
		sval, _ := s.context.ToValue(status)
		rval, _ := s.context.ToValue(res)
		fn.Call(fn, sval, rval)
	}()

	return otto.TrueValue()
}

func (s *Script) scriptHttpshorten(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) < 2 || !call.Argument(0).IsString() || !call.Argument(1).IsFunction() {
		return otto.FalseValue()
	}

	url, fn := call.Argument(0).String(), call.Argument(1)
	go func() {
		res, status := http.Shorten(url)
		sval, _ := s.context.ToValue(status)
		rval, _ := s.context.ToValue(res)
		fn.Call(fn, sval, rval)
	}()

	return otto.TrueValue()
}

func (s *Script) scriptHTMLNew(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) < 1 || !call.ArgumentList[0].IsString() {
		return otto.NullValue()
	}

	url := call.Argument(0).String()
	parser, err := html.NewParser(url, s.context)
	if err != nil {
		log.Errorf("(%s) can't get '%s': %s\n", s.Name(), url, err)
		return otto.NullValue()
	}

	val, err := s.context.ToValue(parser)
	if err != nil {
		log.Errorf("(%s) can't convert parser: %s\n", s.Name(), err)
		return otto.NullValue()
	}

	return val
}
