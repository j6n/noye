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
			matches := findMatches(msg.Text, re)
			if !re.MatchString(msg.Text) {
				continue
			}

			res, err := script.context.ToValue(matches)
			if err != nil {
				continue
			}

			go safeRun(fn, script.Name, val, res)
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
			go safeRun(cmd, script.Name, val)
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

	// add the log method
	if err := ctx.Set("log", func(call otto.FunctionCall) otto.Value {
		// this converts an otto.Value to an interface, so fmt.Sprintf can be used
		toInterface := func(vals []otto.Value) (out []interface{}) {
			for _, val := range vals {
				if res, err := val.Export(); err == nil {
					out = append(out, res)
				}
			}
			return
		}

		// if we got a string
		if call.ArgumentList[0].IsString() {
			// then convert it to a string
			msg := call.Argument(0).String()
			// and check to see if its a Printf-style call, or just a Println
			if len(call.ArgumentList) > 1 {
				msg = fmt.Sprintf(msg, toInterface(call.ArgumentList[1:])...)
			}
			log.Infof("(%s) %s\n", name, msg)
			return otto.TrueValue()
		}
		return otto.FalseValue()
	}); err != nil {
		log.Errorf("(%s) setting log: %s", name, err)
		return err
	}

	// this is kind of a mess, but it returnsa function
	// it can represent a respond or listen object, depending on the path
	build := func(path string) func(otto.FunctionCall) otto.Value {
		return func(call otto.FunctionCall) otto.Value {
			// we require atleast a string and a function for respond and listen
			if len(call.ArgumentList) < 2 || !call.ArgumentList[0].IsString() || !call.ArgumentList[1].IsFunction() {
				return otto.FalseValue()
			}

			// get the input and the callback
			input, fn := call.ArgumentList[0].String(), call.ArgumentList[1]

			// wrap the callback so we can log errors
			wrap := func(env otto.Value, res ...otto.Value) {
				if _, err := fn.Call(fn, append([]otto.Value{env}, res...)); err != nil {
					log.Errorf("(%s,%s,%s) calling fn: %s\n", name, path, input, err)
				}
			}

			switch path {
			// if its a respond, then we'll be using a regex
			case "respond":
				re, err := regexp.Compile(input)
				if err != nil {
					log.Errorf("(%s,%s,%s) compiling re: %s\n", name, path, input, err)
					return otto.FalseValue()
				}
				script.commands[re] = wrap

			// otherwise we just add the callback to the event
			case "listen":
				script.callbacks[input] = append(script.callbacks[input], wrap)
			}
			return otto.TrueValue()
		}
	}

	// add the respond method
	if err := ctx.Set("respond", build("respond")); err != nil {
		log.Errorf("(%s) setting respond: %s\n", name, err)
		return err
	}

	// add the listen method
	if err := ctx.Set("listen", build("listen")); err != nil {
		log.Errorf("(%s) setting listen: %s\n", name, err)
		return err
	}

	// run the actual script
	if _, err := ctx.Run(source); err != nil {
		log.Errorf("(%s) loading script: %s\n", name, err)
		return err
	}

	// if we've gotten this far, the script is valid
	m.scripts[name] = script
	return nil
}

func safeRun(fn scriptFunc, name string, vals ...otto.Value) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("(%s) recovered: %s\n", name, err)
		}
	}()

	fn(vals[0], vals[1:]...)
}

func findMatches(s string, re *regexp.Regexp) map[string]string {
	captures := make(map[string]string)

	match := re.FindStringSubmatch(s)
	if match == nil {
		return captures
	}

	for i, name := range re.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}

		captures[name] = match[i]
	}

	return captures
}
