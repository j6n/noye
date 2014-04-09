package ext

import (
	"io/ioutil"

	"github.com/robertkrimen/otto"
)

var base string

// ReloadBase reloads the base.js script
func (m *Manager) ReloadBase() error {
	data, err := ioutil.ReadFile("base.js")
	if err != nil {
		return err
	}

	log.Infof("loaded base.js script\n")
	base = string(data)
	return nil
}

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

		"_http_get":     script.scriptHttpget,
		"_http_follow":  script.scriptHttpfollow,
		"_http_shorten": script.scriptHttpshorten,
		"_html_new":     script.scriptHTMLNew,
	}

	for k, v := range binding {
		if err := script.context.Set(k, v); err != nil {
			log.Errorf("Couldn't set %s: %s\n", k, err)
			return
		}
	}

	if _, err := script.context.Run(base); err != nil {
		log.Errorf("Couldn't run internal base.js script: %s\n", err)
	}
}
