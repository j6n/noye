package ext

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/j6n/logger"
	"github.com/j6n/noye/noye"
	"github.com/robertkrimen/otto"
)

var log *logger.Logger

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

func New(ctx noye.Bot, logger *logger.Logger) *Manager {
	if log == nil {
		log = logger
	}
	return &Manager{make(map[string]*Script), NewProxyBot(ctx)}
}

func (m *Manager) Respond(msg noye.Message) {
	fields := logger.Fields{
		"manager": "respond",
		"data":    msg,
	}

	for _, script := range m.scripts {
		f := copyFields(fields, logger.Fields{"script": script.Name})

		val, err := script.context.ToValue(msg)
		if err != nil {
			log.WithFields(f).Error(err)
			return
		}

		log.WithFields(f).Debug("attempting")
		for re, fn := range script.commands {
			if !re.MatchString(msg.Text) {
				continue
			}

			go func(val otto.Value, fn scriptFunc) {
				defer func() { recover() }()
				log.WithFields(f).Debug("match, calling")
				fn(val)
			}(val, fn)
		}
	}
}

func (m *Manager) Listen(msg noye.IrcMessage) {
	fields := logger.Fields{
		"manager": "listen",
		"data":    msg,
	}

	for _, script := range m.scripts {
		f := copyFields(fields, logger.Fields{"script": script.Name})

		val, err := script.context.ToValue(msg)
		if err != nil {
			log.WithFields(f).Error(err)
			return
		}

		log.WithFields(f).Debug("attempting")
		cmds, ok := script.callbacks[msg.Command]
		if !ok {
			continue
		}

		log.WithFields(f).Debug("found callbacks")
		for _, cmd := range cmds {
			go func(val otto.Value, fn scriptFunc) {
				defer func() { recover() }()
				log.WithFields(f).Debug("match, calling")
				cmd(val)
			}(val, cmd)
		}
	}
}

func (m *Manager) Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.WithField("plugin", path).Error(err)
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

	fields := logger.Fields{
		"script": name,
		"path":   path,
	}

	log.WithFields(fields).Debug("adding logger")
	ctx.Set("log", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) == 1 && call.ArgumentList[0].IsString() {
			log.WithFields(copyFields(fields, logger.Fields{"event": "log"})).Info(call.ArgumentList[0].String())
			return otto.TrueValue()
		}
		return otto.FalseValue()
	})

	fields = copyFields(fields, logger.Fields{"event": "load"})
	log.WithFields(fields).Debug("adding respond")
	ctx.Set("respond", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 2 || !call.ArgumentList[0].IsString() || !call.ArgumentList[1].IsFunction() {
			return otto.FalseValue()
		}

		str, fn := call.ArgumentList[0].String(), call.ArgumentList[1]
		wrap := func(env otto.Value) {
			if _, err := fn.Call(fn, env); err != nil {
				log.WithFields(copyFields(fields, logger.Fields{"call": "cmd", "subject": "respond"})).Error(err)
			}
		}

		re, err := regexp.Compile(str)
		if err != nil {
			log.WithFields(copyFields(fields, logger.Fields{"compile": "regex", "string": str})).Error(err)
			return otto.FalseValue()
		}

		script.commands[re] = wrap
		return otto.TrueValue()
	})

	log.WithFields(fields).Debug("adding listen")
	ctx.Set("listen", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 2 || !call.ArgumentList[0].IsString() || !call.ArgumentList[1].IsFunction() {
			return otto.FalseValue()
		}

		cmd, fn := call.ArgumentList[0].String(), call.ArgumentList[1]
		wrap := func(env otto.Value) {
			if _, err := fn.Call(fn, env); err != nil {
				log.WithFields(copyFields(fields, logger.Fields{"call": "cmd", "subject": "listen"})).Error(err)
			}
		}
		script.callbacks[cmd] = append(script.callbacks[cmd], wrap)
		return otto.TrueValue()
	})

	log.WithFields(fields).Debug("running script")
	_, err := ctx.Run(source)
	if err != nil {
		log.WithFields(fields).Error(err)
		return err
	}

	m.scripts[name] = script
	return nil
}

func copyFields(origin, input logger.Fields) logger.Fields {
	out := logger.Fields{}
	for k, v := range origin {
		out[k] = v
	}

	for k, v := range input {
		out[k] = v
	}

	return out
}
