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
	"http": function(url) {	return new HttpClient(url); },
};

function HttpClient(url) {
	this._url = url;	
}

HttpClient.prototype.get = function() {	
	return _httpclient_get(this._url);
}
`

func (m *Manager) defaults(vm *otto.Otto) {
	if err := vm.Set("_core_reply", m.proxy.Reply); err != nil {
		// TODO log error
		err = nil
		return
	}

	if err := vm.Set("_core_bot", m.proxy); err != nil {
		// TODO log error
		err = nil
		return
	}

	if err := vm.Set("_core_load", m.Load); err != nil {
		// TODO log error
		err = nil
		return
	}

	if err := vm.Set("_httpclient_get", http_get); err != nil {
		// TODO log error
		err = nil
		return
	}

	if _, err := vm.Run(base); err != nil {
		// TODO log error
		err = nil
	}
}

func http_post(args ...string) string {
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

	defer resp.Body.Close()
	buf := new(bytes.Buffer)

	if _, err := io.Copy(buf, resp.Body); err != nil {
		return ""
	}

	return buf.String()
}

func http_get(args ...string) string {
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

	defer resp.Body.Close()
	buf := new(bytes.Buffer)

	if _, err := io.Copy(buf, resp.Body); err != nil {
		return ""
	}

	return buf.String()
}
