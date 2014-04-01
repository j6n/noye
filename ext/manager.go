package ext

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"github.com/j6n/noye/noye"
	"github.com/robertkrimen/otto"
)

type Script struct {
	Name, Path, Source string

	commands  map[*regexp.Regexp]scriptFunc
	callbacks map[string][]scriptFunc

	context *otto.Otto
}

type scriptFunc func(otto.Value)

type Manager struct {
	scripts map[string]*Script
	proxy   *ProxyBot
}

func New(ctx noye.Bot) *Manager {
	return &Manager{make(map[string]*Script), NewProxyBot(ctx)}
}

func (m *Manager) Respond(msg noye.Message) {
	for _, script := range m.scripts {
		val, err := script.context.ToValue(msg)
		if err != nil {
			return
		}

		for re, fn := range script.commands {
			if !re.MatchString(msg.Text) {
				continue
			}

			go func(val otto.Value, fn scriptFunc) {
				defer func() { recover() }()
				fn(val)
			}(val, fn)
		}
	}
}

func (m *Manager) Listen(msg noye.IrcMessage) {
	for _, script := range m.scripts {
		val, err := script.context.ToValue(msg)
		if err != nil {
			return
		}

		cmds, ok := script.callbacks[msg.Command]
		if !ok {
			continue
		}

		for _, cmd := range cmds {
			go func(val otto.Value, fn scriptFunc) {
				defer func() { recover() }()
				cmd(val)
			}(val, cmd)
		}
	}
}

func (m *Manager) Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return m.load(string(data), path)
}

func (m *Manager) Reload(name string) error {
	if script, ok := m.scripts[name]; ok {
		delete(m.scripts, name)
		return m.load(script.Source, script.Path)
	}

	// script not loaded
	return fmt.Errorf("%s is not loaded", name)
}

func (m *Manager) load(source, path string) error {
	name := filepath.Base(path)
	ctx := otto.New()

	script := &Script{
		Name:   name,
		Path:   path,
		Source: source,

		commands:  make(map[*regexp.Regexp]scriptFunc),
		callbacks: make(map[string][]scriptFunc),

		context: ctx,
	}

	// init proxy bot
	m.defaults(ctx)

	ctx.Set("log", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) == 1 && call.ArgumentList[0].IsString() {
			// TODO log bot stuff here
			return otto.TrueValue()
		}
		return otto.FalseValue()
	})

	ctx.Set("respond", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 2 || !call.ArgumentList[0].IsString() || !call.ArgumentList[1].IsFunction() {
			return otto.FalseValue()
		}

		str, fn := call.ArgumentList[0].String(), call.ArgumentList[1]
		wrap := func(env otto.Value) {
			if _, err := fn.Call(fn, env); err != nil {
				// TODO log error
				err = nil
			}
		}

		re, err := regexp.Compile(str)
		if err != nil {
			return otto.FalseValue()
		}

		script.commands[re] = wrap
		return otto.TrueValue()
	})

	ctx.Set("listen", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 2 || !call.ArgumentList[0].IsString() || !call.ArgumentList[1].IsFunction() {
			return otto.FalseValue()
		}

		cmd, fn := call.ArgumentList[0].String(), call.ArgumentList[1]
		wrap := func(env otto.Value) {
			if _, err := fn.Call(fn, env); err != nil {
				// TODO log error
				err = nil
			}
		}
		script.callbacks[cmd] = append(script.callbacks[cmd], wrap)
		return otto.TrueValue()
	})

	if _, err := ctx.Run(source); err != nil {
		return err
	}

	m.scripts[name] = script
	return nil
}
