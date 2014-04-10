package script

import (
	"io/ioutil"

	"github.com/j6n/noye/core/lib"
	"github.com/j6n/noye/core/logger"
	"github.com/j6n/noye/core/store"

	"github.com/robertkrimen/otto"
)

var log = logger.Get()
var base string

// ReloadBase reloads the base.js script
func (m *Manager) ReloadBase() (err error) {
	var data []byte
	if data, err = ioutil.ReadFile("base.js"); err == nil {
		base = string(data)
	}
	return
}

// set the defaults for the script
func (m *Manager) setDefaults(s *Script) (err error) {
	binding := map[string]interface{}{
		"_noye_bot": m.context,

		"_core_manager": m,
		"_core_scripts": getScripts(m),

		"_core_storage_load": getMethod(s),
		"_core_storage_save": setMethod(s),

		"_share_init":   initMethod(s),
		"_share_sub":    subMethod(s),
		"_share_unsub":  unsubMethod(s),
		"_share_update": updateMethod(s),

		"_http_get":     httpGetMethod(s),
		"_http_follow":  httpFollowMethod(s),
		"_http_shorten": shortenMethod(s),
		"_html_new":     newParserMethod(s),
	}

	for k, v := range binding {
		if err = s.context.Set(k, v); err != nil {
			return
		}
	}

	if _, err = s.context.Run(base); err != nil {
		log.Errorf("unable to run base.js for '%s': %s", s.Name(), err)
	}

	return
}

// gets a listing of scripts
func getScripts(m *Manager) func(otto.FunctionCall) otto.Value {
	fn := func(call otto.FunctionCall) otto.Value {
		var resp = struct {
			Names   []string
			Details map[string]string
		}{make([]string, 0), make(map[string]string)}

		for k, v := range m.scripts {
			resp.Names = append(resp.Names, k)
			resp.Details[v.Name()] = v.Path()
		}

		val, err := m.vm.ToValue(resp)
		if err != nil {
			return otto.UndefinedValue()
		}

		return val
	}

	return fn
}

// stores a key:val
func setMethod(s *Script) func(otto.FunctionCall) otto.Value {
	fn := func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 2 || !call.Argument(0).IsString() || !call.Argument(1).IsString() {
			return otto.FalseValue()
		}

		key, data := call.Argument(0).String(), call.Argument(1).String()
		if err := store.Set(s.Name(), key, data); err != nil {
			log.Errorf("(%s) setting val at '%s': %s", s.Name(), key, err)
			return otto.FalseValue()
		}

		return otto.TrueValue()
	}

	return fn
}

// gets a val for a key
func getMethod(s *Script) func(otto.FunctionCall) otto.Value {
	fn := func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 || !call.Argument(0).IsString() {
			return otto.UndefinedValue()
		}

		key := call.Argument(0).String()
		data, err := store.Get(s.Name(), key)
		if err != nil {
			log.Errorf("(%s) getting val at '%s': %s", s.Name(), key, err)
			return otto.UndefinedValue()
		}

		val, _ := s.context.ToValue(data)
		return val
	}

	return fn
}

func initMethod(s *Script) func(otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 2 || !call.Argument(0).IsString() || !call.Argument(1).IsFunction() {
			return otto.FalseValue()
		}

		key, fn := call.Argument(0).String(), call.Argument(1).Object()
		s.inits[key] = append(s.inits[key], fn)
		log.Debugf("(%s) added init '%s'", s.Name(), key)
		return otto.TrueValue()
	}
}

// subscribes to the message queue
func subMethod(s *Script) func(otto.FunctionCall) otto.Value {
	fn := func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 2 || !call.Argument(0).IsString() || !call.Argument(1).IsFunction() {
			return otto.FalseValue()
		}

		key, fn := call.Argument(0).String(), call.Argument(1).Object()
		s.updates[key] = append(s.updates[key], fn)
		log.Debugf("(%s) added update '%s'", s.Name(), key)
		return otto.TrueValue()
	}

	return fn
}

// unsubscribes from the message queue
func unsubMethod(s *Script) func(otto.FunctionCall) otto.Value {
	fn := func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 1 || !call.Argument(0).IsNumber() {
			return otto.FalseValue()
		}

		id, err := call.Argument(0).ToInteger()
		if err != nil {
			log.Errorf("(%s) wasn't given an mq id: %s", s.Name(), err)
			return otto.FalseValue()
		}

		mq.Unsubscribe(id)
		for i, sub := range s.subs {
			if sub == id {
				s.subs = s.subs[:i+copy(s.subs[i:], s.subs[i+1:])]
			}
		}
		return otto.TrueValue()
	}

	return fn
}

// broadcasts a message for a key
func updateMethod(s *Script) func(otto.FunctionCall) otto.Value {
	fn := func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 2 || !call.Argument(0).IsString() || !call.Argument(1).IsString() {
			return otto.FalseValue()
		}

		key, val := call.Argument(0).String(), call.Argument(1).String()
		mq.Update(key, val, false)
		return otto.TrueValue()
	}

	return fn
}

// http get for a string
func httpGetMethod(s *Script) func(otto.FunctionCall) otto.Value {
	fn := func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 1 || !call.Argument(0).IsString() {
			return otto.NullValue()
		}

		headers := make(map[string]string)
		if len(call.ArgumentList) == 2 {
			if obj, err := call.Argument(1).Export(); err == nil {
				if m, ok := obj.(map[string]interface{}); ok {
					for k, v := range m {
						headers[k] = v.(string)
					}
				}
			}
		}

		res, status := lib.Get(call.Argument(0).String(), headers)
		sval, _ := s.context.ToValue(status)
		rval, _ := s.context.ToValue(res)
		data, _ := s.context.ToValue(map[string]interface{}{"status": sval, "body": rval})
		return data
	}
	return fn
}

// http follow a string
func httpFollowMethod(s *Script) func(otto.FunctionCall) otto.Value {
	fn := func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 1 || !call.Argument(0).IsString() {
			return otto.NullValue()
		}

		res, status := lib.Follow(call.Argument(0).String())
		sval, _ := s.context.ToValue(status)
		rval, _ := s.context.ToValue(res)
		data, _ := s.context.ToValue(map[string]interface{}{"status": sval, "body": rval})
		return data
	}
	return fn
}

// returns a goo.gl url for a string
func shortenMethod(s *Script) func(otto.FunctionCall) otto.Value {
	fn := func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 2 || !call.Argument(0).IsString() {
			return otto.NullValue()
		}

		res, status := lib.Shorten(call.Argument(0).String())
		sval, _ := s.context.ToValue(status)
		rval, _ := s.context.ToValue(res)

		data, _ := s.context.ToValue(map[string]interface{}{"status": sval, "body": rval})
		return data
	}

	return fn
}

// returns a new httml parser
func newParserMethod(s *Script) func(otto.FunctionCall) otto.Value {
	fn := func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 1 || !call.Argument(0).IsString() {
			return otto.NullValue()
		}

		url := call.Argument(0).String()
		parser, err := lib.NewParser(url, s.context)
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

	return fn
}

func listen(s *Script, obj *otto.Object, ch chan string) {
	fn := obj.Value()

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
}
