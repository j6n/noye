package ext

import (
	"github.com/j6n/noye/store"

	"github.com/robertkrimen/otto"
)

const base = `
noye = _noye_bot;
core = {
	"manager": _core_manager,
	"scripts": _core_scripts,
	"load": _core_storage_load,
	"store": _core_storage_save,
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

	binding := map[string]interface{}{
		"_noye_bot":          m.context,
		"_core_manager":      m,
		"_core_scripts":      scripts,
		"_core_storage_load": get,
		"_core_storage_save": set,
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
