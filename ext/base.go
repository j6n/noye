package ext

import "github.com/robertkrimen/otto"

const base = `
noye = _noye_bot;

core = {
	"manager": _core_manager,
	"scripts": _core_scripts,
	"load":    _core_storage_load,
	"save":    _core_storage_save,
};

share = {
	"init":        _share_init,
	"update":      _share_update,
	"subscribe":   _share_sub,
	"unsubscribe": _share_unsub,
};

http = {
	"get": _http_get,
};
`

func (m *Manager) setDefaults(script *Script) {
	getScriptsFor := func() otto.Value {
		var resp = struct {
			Scripts []string
			Details map[string]string
		}{make([]string, 0), make(map[string]string)}

		for _, s := range m.Scripts() {
			resp.Scripts = append(resp.Scripts, s.Name())
			resp.Details[s.Name()] = s.Path()
		}

		val, err := script.context.ToValue(resp)
		if err != nil {
			return otto.UndefinedValue()
		}

		return val
	}

	binding := map[string]interface{}{
		"_noye_bot": m.context,

		"_core_manager":      m,
		"_core_scripts":      getScriptsFor,
		"_core_storage_load": script.scriptGet,
		"_core_storage_save": script.scriptSet,

		"_share_init":   script.scriptSubInit(true),
		"_share_sub":    script.scriptSubInit(false),
		"_share_unsub":  script.scriptUnsub,
		"_share_update": script.scriptUpdate,

		"_http_get": script.scriptHttpget,
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
