package ext

import (
	"github.com/j6n/noye/http"
	"github.com/j6n/noye/store"

	"github.com/robertkrimen/otto"
)

const base = `
noye = _noye_bot;
core = {
	"manager": _core_manager,
	"scripts": _core_scripts,
	"load":    _core_storage_load,
	"save":    _core_storage_save,
};
share = {
	"update":      _share_update,
	"subscribe":   _share_sub,
	"unsubscribe": _share_unsub,
};

http = {
	"get": _http_get,
};
`

func (m *Manager) setDefaults(script *Script) {
	binding := map[string]interface{}{
		"_noye_bot": m.context,

		"_core_manager":      m,
		"_core_scripts":      getScripts(m, script),
		"_core_storage_load": scriptGet(script),
		"_core_storage_save": scriptSet(script),

		"_share_sub":    scriptSub(script),
		"_share_unsub":  scriptUnsub(script),
		"_share_update": scriptUpdate(script),

		"_http_get": scriptHttpget(script),
	}

	for k, v := range binding {
		if err := script.context.Set(k, v); err != nil {
			log.Errorf("Couldn't set %s: %s\n", k, err)
			return
		}
	}

	if _, err := script.context.Run(base); err != nil {
		log.Errorf("Couldn't run base script: %s\n", err)
	}
}

var mq = store.NewQueue()

func getScripts(m *Manager, s *Script) otto.Value {
	var resp = struct {
		Scripts []string
		Details map[string]string
	}{make([]string, 0), make(map[string]string)}

	for k, v := range m.scripts {
		resp.Scripts = append(resp.Scripts, k)
		resp.Details[k] = v.Path()
	}

	val, err := s.context.ToValue(resp)
	if err != nil {
		return otto.UndefinedValue()
	}

	return val
}

func scriptSet(s *Script) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
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
}

func scriptGet(s *Script) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
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
}

func scriptSub(s *Script) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 2 || !call.ArgumentList[0].IsString() || !call.ArgumentList[1].IsFunction() {
			return otto.NullValue()
		}

		key, fn := call.ArgumentList[0].String(), call.ArgumentList[1]
		id, ch := mq.Subscribe(key)

		val, err := s.context.ToValue(id)
		if err != nil {
			log.Errorf("(%s) convert val (sub): %s", s.Name(), err)
			return otto.NullValue()
		}

		go func() {
			for data := range ch {
				val, err := s.context.ToValue(data)
				if err != nil {
					log.Errorf("(%s) convert val '%s': %s", s.Name(), data, err)
					val = otto.NullValue()
				}
				fn.Call(fn, val)
			}
		}()

		s.subs = append(s.subs, id)
		return val
	}
}

func scriptUnsub(s *Script) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
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
}

func scriptUpdate(s *Script) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 2 || !call.ArgumentList[0].IsString() || !call.ArgumentList[1].IsString() {
			return otto.FalseValue()
		}

		key, val := call.ArgumentList[0].String(), call.ArgumentList[1].String()
		mq.Update(key, val)
		return otto.TrueValue()
	}
}

func scriptHttpget(s *Script) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 2 || !call.ArgumentList[0].IsString() || !call.ArgumentList[1].IsFunction() {
			return otto.FalseValue()
		}

		var headers map[string]string
		if len(call.ArgumentList) > 2 {
			obj, err := call.ArgumentList[3].Export()
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
			status, res := http.Get(url, headers)
			sval, _ := s.context.ToValue(status)
			rval, _ := s.context.ToValue(res)

			fn.Call(fn, sval, rval)
		}()
		return otto.TrueValue()
	}
}
