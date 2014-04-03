package ext

import "github.com/robertkrimen/otto"

const base = `
noye = _noye_bot;
core = {
	"manager": _core_manager,
	"scripts": _core_scripts,
};
`

func (m *Manager) setDefaults(vm *otto.Otto) {
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

	binding := map[string]interface{}{
		"_noye_bot":     m.context,
		"_core_manager": m,
		"_core_scripts": scripts,
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
