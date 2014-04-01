package ext

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"github.com/robertkrimen/otto"
)

const base = `
noye = {
	"reply": function() { _core_reply.apply(null, arguments); },
	"bot": _core_bot,
};

core = {
	"load": function() { _core_load.apply(null, arguments); },
	"http": function(url) {	return new _http(url); },
};

function _http(url) {
	this._url = url;	
}

_http.prototype.get = function() {	
	return _http_get(this._url);
}
`

func (m *Manager) defaults(vm *otto.Otto) {
	if err := vm.Set("_core_reply", m.proxy.Reply); err != nil {
		log.Errorf("Couldn't set _core_reply: %s\n", err)
		return
	}

	if err := vm.Set("_core_bot", m.proxy); err != nil {
		log.Errorf("Couldn't set _core_bot: %s\n", err)
		return
	}

	if err := vm.Set("_core_load", m.Load); err != nil {
		log.Errorf("Couldn't set _core_load: %s\n", err)
		return
	}

	if err := vm.Set("_http_get", httpGet); err != nil {
		log.Errorf("Couldn't set _http_get: %s\n", err)
		return
	}

	if _, err := vm.Run(base); err != nil {
		log.Errorf("Couldn't run base script: %s\n", err)
	}
}

func httpPost(args ...string) string {
	url := strings.Trim(args[0], `"`)
	client := &http.Client{}

	var body *bytes.Buffer
	if len(args) > 2 {
		body = bytes.NewBufferString(strings.Trim(args[1], `"`)[1:])
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return ""
	}

	if len(args) > 2 {
		var headers map[string]string
		err = json.Unmarshal([]byte(args[2]), &headers)
		if err != nil {
			return ""
		}

		for key, value := range headers {
			req.Header.Add(key, value)
		}
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
