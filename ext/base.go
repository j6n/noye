package ext

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/robertkrimen/otto"
)

const base = `
noye = _noye_bot;
core = {
	"http": function(url) {	return new _http(url); },
	"scripts": _core_scripts,
};

function _http(url) {
	this._url = url;	
}

_http.prototype.get = function() {	
	return _http_get(this._url);
}
`

func (m *Manager) defaults(vm *otto.Otto) {
	set := func(name string, what interface{}) bool {
		if err := vm.Set(name, what); err != nil {
			log.Errorf("Couldn't set %s: %s\n", name, err)
			return false
		}

		return true
	}

	if !set("_noye_bot", m.context) {
		return
	}

	if !set("_core_scripts", func() otto.Value {
		var resp = struct {
			Scripts []string
			Details map[string]string
		}{make([]string, 0), make(map[string]string)}

		for k, v := range m.scripts {
			resp.Scripts = append(resp.Scripts, k)
			resp.Details[k] = v.Path
		}

		val, err := vm.ToValue(resp)
		if err != nil {
			return otto.UndefinedValue()
		}

		return val
	}) {
		return
	}

	if !set("_http_get", httpGet) {
		return
	}

	if _, err := vm.Run(base); err != nil {
		log.Errorf("Couldn't run base script: %s\n", err)
	}
}

func httpGet(args ...string) string {
	url := strings.Trim(args[0], `"`)
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			// do nothing
		}
	}()
	buf := new(bytes.Buffer)

	if _, err := io.Copy(buf, resp.Body); err != nil {
		return ""
	}

	return buf.String()
}
