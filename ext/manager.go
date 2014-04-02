package ext

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/j6n/noye/logger"
	"github.com/j6n/noye/noye"
	"github.com/robertkrimen/otto"
)

var log = logger.Get()

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

			go safeRun(val, fn, script.Name)
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
			go safeRun(val, cmd, script.Name)
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

// Eval evaluates the source, returning any errors
func (m *Manager) Eval(source string) error {
	return m.load(source, "/Test.js")
}

func (m *Manager) load(source, path string) error {
	name := filepath.Base(path)
	script := newScript(name, path, source)

	// copy pointer
	ctx := script.context

	// init proxy bot
	m.defaults(ctx)

	if err := ctx.Set("log", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) == 1 && call.ArgumentList[0].IsString() {
			log.Infof("(%s) %s\n", name, call.Argument(0).String())
			return otto.TrueValue()
		}
		return otto.FalseValue()
	}); err != nil {
		log.Errorf("(%s) setting log: %s", name, err)
		return err
	}

	build := func(path string) func(otto.FunctionCall) otto.Value {
		return func(call otto.FunctionCall) otto.Value {
			if len(call.ArgumentList) < 2 || !call.ArgumentList[0].IsString() || !call.ArgumentList[1].IsFunction() {
				return otto.FalseValue()
			}

			input, fn := call.ArgumentList[0].String(), call.ArgumentList[1]
			wrap := func(env otto.Value) {
				if _, err := fn.Call(fn, env); err != nil {
					log.Errorf("(%s,%s,%s) calling fn: %s\n", name, path, input, err)
				}
			}

			switch path {
			case "respond":
				re, err := regexp.Compile(input)
				if err != nil {
					log.Errorf("(%s,%s,%s) compiling re: %s\n", name, path, input, err)
					return otto.FalseValue()
				}

				script.commands[re] = wrap
			case "listen":
				script.callbacks[input] = append(script.callbacks[input], wrap)
			}
			return otto.TrueValue()
		}
	}

	if err := ctx.Set("respond", build("respond")); err != nil {
		log.Errorf("(%s) setting respond: %s\n", name, err)
		return err
	}

	if err := ctx.Set("listen", build("listen")); err != nil {
		log.Errorf("(%s) setting listen: %s\n", name, err)
		return err
	}

	if _, err := ctx.Run(source); err != nil {
		log.Errorf("(%s) loading script: %s\n", name, err)
		return err
	}

	m.scripts[name] = script
	return nil
}

func safeRun(val otto.Value, fn scriptFunc, name string) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("(%s) recovered: %s\n", name, err)
		}
	}()

	fn(val)
}
