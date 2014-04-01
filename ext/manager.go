package ext

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/j6n/noye/noye"
	"github.com/robertkrimen/otto"
)

// Manager holds a bunch of scripts and a safe proxy to the bot
type Manager struct {
	scripts map[string]*Script
	proxy   *ProxyBot
}

// New returns a new Manager
func New(ctx noye.Bot) *Manager {
	return &Manager{make(map[string]*Script), NewProxyBot(ctx)}
}

// Respond takes a noye.Message and delegates it to the scripts
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
				defer func() { _ = recover() }()
				fn(val)
			}(val, fn)
		}
	}
}

// Listen takes a noye.IrcMessage and delegates it to the scripts
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
				defer func() { _ = recover() }()
				cmd(val)
			}(val, cmd)
		}
	}
}

// Load tries to load the file located at the path
func (m *Manager) Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return m.load(string(data), path)
}

// Reload tries to reload the named script
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
	script := newScript(name, path, source)

	// copy pointer
	ctx := script.context

	// init proxy bot
	m.defaults(ctx)

	_ = ctx.Set("log", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) == 1 && call.ArgumentList[0].IsString() {
			// TODO log bot stuff here
			return otto.TrueValue()
		}
		return otto.FalseValue()
	})

	build := func(path string) func(otto.FunctionCall) otto.Value {
		return func(call otto.FunctionCall) otto.Value {
			if len(call.ArgumentList) < 2 || !call.ArgumentList[0].IsString() || !call.ArgumentList[1].IsFunction() {
				return otto.FalseValue()
			}

			input, fn := call.ArgumentList[0].String(), call.ArgumentList[1]
			wrap := func(env otto.Value) {
				if _, err := fn.Call(fn, env); err != nil {
					// TODO log error
					err = nil
				}
			}

			switch path {
			case "respond":
				re, err := regexp.Compile(input)
				if err != nil {
					return otto.FalseValue()
				}

				script.commands[re] = wrap
			case "listen":
				script.callbacks[input] = append(script.callbacks[input], wrap)
			}
			return otto.TrueValue()
		}
	}

	_ = ctx.Set("respond", build("respond"))
	_ = ctx.Set("listen", build("listen"))

	if _, err := ctx.Run(source); err != nil {
		return err
	}

	m.scripts[name] = script
	return nil
}
