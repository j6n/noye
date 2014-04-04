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

func (m *Manager) setDefaults(vm *otto.Otto, script *Script) {
	scripts := func() otto.Value {
		var resp = struct {
			Scripts []string
			Details map[string]string
		}{make([]string, 0), make(map[string]string)}

		for k, v := range m.scripts {
			resp.Scripts = append(resp.Scripts, k)
			resp.Details[k] = v.Path()
		}

		val, err := vm.ToValue(resp)
		if err != nil {
			return otto.UndefinedValue()
		}

		return val
	}

	set := func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 2 || !call.ArgumentList[0].IsString() || !call.ArgumentList[1].IsString() {
			return otto.FalseValue()
		}

		key, data := call.ArgumentList[0].String(), call.ArgumentList[1].String()
		if err := store.Set(script.Name(), key, data); err != nil {
			log.Errorf("(%s) setting val at '%s': %s", script.Name(), key, err)
			return otto.FalseValue()
		}

		return otto.TrueValue()
	}

	get := func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 || !call.ArgumentList[0].IsString() {
			return otto.UndefinedValue()
		}

		key := call.ArgumentList[0].String()
		data, err := store.Get(script.Name(), key)
		if err != nil {
			log.Errorf("(%s) getting val at '%s': %s", script.Name(), key, err)
			return otto.UndefinedValue()
		}

		val, _ := vm.ToValue(data)
		return val
	}

	sub := func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 2 || !call.ArgumentList[0].IsString() || !call.ArgumentList[1].IsFunction() {
			return otto.NullValue()
		}

		key, fn := call.ArgumentList[0].String(), call.ArgumentList[1]
		id, ch := mq.Subscribe(key)

		val, err := script.context.ToValue(id)
		if err != nil {
			log.Errorf("(%s) convert val (sub): %s", script.Name(), err)
			return otto.NullValue()
		}

		go func() {
			for data := range ch {
				val, err := script.context.ToValue(data)
				if err != nil {
					log.Errorf("(%s) convert val '%s': %s", script.Name(), data, err)
					val = otto.NullValue()
				}
				fn.Call(fn, val)
			}
		}()

		script.subs = append(script.subs, id)
		return val
	}

	unsub := func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 1 || !call.ArgumentList[0].IsNumber() {
			return otto.FalseValue()
		}
		id, err := call.ArgumentList[0].ToInteger()
		if err != nil {
			log.Errorf("(%s) wasn't given an unsub id: %s", script.Name(), err)
			return otto.TrueValue()
		}

		mq.Unsubscribe(id)
		for i, s := range script.subs {
			if s == id {
				script.subs = script.subs[:i+copy(script.subs[i:], script.subs[i+1:])]
			}
		}
		return otto.TrueValue()
	}

	update := func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 2 || !call.ArgumentList[0].IsString() || !call.ArgumentList[1].IsString() {
			return otto.FalseValue()
		}

		key, val := call.ArgumentList[0].String(), call.ArgumentList[1].String()
		mq.Update(key, val)
		return otto.TrueValue()
	}

	httpget := func(call otto.FunctionCall) otto.Value {
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
			sval, _ := script.context.ToValue(status)
			rval, _ := script.context.ToValue(res)

			fn.Call(fn, sval, rval)
		}()
		return otto.TrueValue()
	}

	binding := map[string]interface{}{
		"_noye_bot": m.context,

		"_core_manager":      m,
		"_core_scripts":      scripts,
		"_core_storage_load": get,
		"_core_storage_save": set,

		"_share_sub":    sub,
		"_share_unsub":  unsub,
		"_share_update": update,

		"_http_get": httpget,
	}

	for k, v := range binding {
		if err := vm.Set(k, v); err != nil {
			log.Errorf("Couldn't set %s: %s\n", k, err)
			return
		}
	}

	if _, err := vm.Run(base); err != nil {
		log.Errorf("Couldn't run base script: %s\n", err)
	}
}

var mq = store.NewQueue()
