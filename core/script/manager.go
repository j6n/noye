package script

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/j6n/noye/core/logger"
	"github.com/j6n/noye/noye"

	"github.com/robertkrimen/otto"
)

var scriptLogger = logger.New()

// Manager holds a bunch of scripts and a safe proxy to the bot
type Manager struct {
	scripts map[string]*Script
	vm      *otto.Otto
	context noye.Bot
}

// New returns a new Manager
func New(ctx noye.Bot) (m *Manager) {
	m = &Manager{
		scripts: make(map[string]*Script),
		vm:      otto.New(),
		context: ctx,
	}
	if err := m.ReloadBase(); err != nil {
		log.Fatalf("unable to reload base: %s", err)
	}
	return
}

// Respond takes a noye.Message and delegates it to the scripts
func (m *Manager) Respond(msg noye.Message) {
	wrapped := m.wrapMessage(msg)
	log.Debugf("(%s:%t) <%s> %s", wrapped.Target, wrapped.Public, wrapped.From, wrapped.Text)
	val, err := m.vm.ToValue(wrapped)
	if err != nil {
		log.Errorf("unable to convert wrapped message: %s", err)
		return
	}

	for _, script := range m.scripts {
		for re, fn := range script.commands {
			if !re.MatchString(msg.Text) {
				continue
			}

			matches := findMatches(msg.Text, re)
			res, err := script.context.ToValue(matches)
			if err != nil {
				continue
			}

			go safeRun(fn, script.Name(), val, res)
		}
	}
}

// Listen takes a noye.IrcMessage and delegates it to the scripts
func (m *Manager) Listen(msg noye.IrcMessage) {
	val, err := m.vm.ToValue(msg)
	if err != nil {
		log.Errorf("unable to convert wrapped message: %s", err)
		return
	}

	for _, script := range m.scripts {
		cmds, ok := script.callbacks[msg.Command]
		if !ok {
			continue
		}

		for _, cmd := range cmds {
			go safeRun(cmd, script.Name(), val)
		}
	}
}

// LoadFile tries to load the file located at the path
func (m *Manager) LoadFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	name := filepath.Base(path)
	if err := m.Reload(name); err == nil {
		log.Infof("reloaded '%s'", name)
		return nil
	}

	return m.load(string(data), name, path)
}

// Reload tries to reload the named script
func (m *Manager) Reload(name string) error {
	if script, ok := m.scripts[name]; ok {
		m.unload(script)
		source, err := ioutil.ReadFile(script.Path())
		if err != nil {
			return err
		}

		return m.load(string(source), name, script.Path())
	}

	return fmt.Errorf("%s is not loaded", name)
}

// Unload tries to remove a named script
func (m *Manager) Unload(name string) (err error) {
	if script, ok := m.scripts[name]; ok {
		m.unload(script)
	}

	return fmt.Errorf("%s is not loaded", name)
}

// UnloadAll tries to unload all of the scripts
func (m *Manager) UnloadAll() {
	for k := range m.scripts {
		m.Unload(k)
	}
}

func (m *Manager) unload(s *Script) {
	for _, sub := range s.subs {
		mq.Unsubscribe(sub)
	}
	s.Cleanup()
	delete(m.scripts, s.Name())
	return
}

func (m *Manager) load(source, name, path string) (err error) {
	script := newScript(name, path, source)
	ctx := script.context

	// init default js methods
	m.setDefaults(script)

	// add the default methods
	if err = checkAnyError(
		ctx.Set("log", logMethod(script)),
		ctx.Set("cleanup", cleanupMethod(script)),
		ctx.Set("respond", respondMethod(script)),
		ctx.Set("listen", listenMethod(script)),
	); err != nil {
		return
	}

	lodash, _ := lodashminjs()
	ctx.Run(string(lodash))

	// run the actual script
	if _, err = ctx.Run(source); err == nil {
		m.scripts[name] = script
		log.Debugf("(%s) ran script", script.Name())
		script.initialize()
		log.Debugf("(%s) initialize ctx", script.Name())
	} else {
		m.unload(script)
	}
	return
}

func logMethod(s *Script) func(otto.FunctionCall) otto.Value {
	fn := func(call otto.FunctionCall) otto.Value {

		if len(call.ArgumentList) == 0 {
			return otto.FalseValue()
		}

		if call.Argument(0).IsString() {
			msg := call.Argument(0).String()
			if len(call.ArgumentList) > 1 {
				msg = fmt.Sprintf(msg, valuesToInterfaces(call.ArgumentList[1:])...)
			}

			scriptLogger.Infof("(%s) %s", s.Name(), msg)
			return otto.TrueValue()
		}

		scriptLogger.Infof("(%s) %s", s.Name(), valuesToString(call.ArgumentList))
		return otto.TrueValue()
	}

	return fn
}

func cleanupMethod(s *Script) func(otto.FunctionCall) otto.Value {
	fn := func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) == 0 || !call.Argument(0).IsFunction() {
			return otto.FalseValue()
		}

		fn := call.Argument(0)
		wrap := func(env otto.Value, res ...otto.Value) {
			if _, err := fn.Call(fn); err != nil {
				log.Errorf("(%s) calling cleanup: %s", s.Name(), err)
			}
		}

		s.cleanup = append(s.cleanup, wrap)
		return otto.TrueValue()
	}

	return fn
}

func respondMethod(s *Script) func(otto.FunctionCall) otto.Value {
	fn := func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 2 || !call.Argument(0).IsString() || !call.Argument(1).IsFunction() {
			return otto.FalseValue()
		}

		input, fn := call.Argument(0).String(), call.Argument(1)
		wrap := func(env otto.Value, res ...otto.Value) {
			vals := []interface{}{env}
			for _, r := range res {
				vals = append(vals, r)
			}

			if _, err := fn.Call(fn, vals...); err != nil {
				log.Errorf("(%s: %s) calling respond func: %s", s.Name(), input, err)
			}
		}

		re, err := regexp.Compile(input)
		if err != nil {
			log.Errorf("(%s: %s) compiling regex: %s", s.Name(), input, err)
			return otto.FalseValue()
		}
		s.commands[re] = wrap
		return otto.TrueValue()
	}

	return fn
}

func listenMethod(s *Script) func(otto.FunctionCall) otto.Value {
	fn := func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 2 || !call.Argument(0).IsString() || !call.Argument(1).IsFunction() {
			return otto.FalseValue()
		}

		input, fn := call.Argument(0).String(), call.Argument(1)
		wrap := func(env otto.Value, res ...otto.Value) {
			vals := []interface{}{env}
			for _, r := range res {
				vals = append(vals, r)
			}

			if _, err := fn.Call(fn, vals...); err != nil {
				log.Errorf("(%s: %s) calling listen func: %s", s.Name(), input, err)
			}
		}

		s.callbacks[input] = append(s.callbacks[input], wrap)
		return otto.TrueValue()
	}

	return fn
}
