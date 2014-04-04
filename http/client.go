package http

import (
	"bytes"
	"io"
	"net/http"
)

func Get(url string, headers map[string]string) (int, string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return http.StatusBadRequest, ""
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return resp.StatusCode, ""
	}

	defer resp.Body.Close()
	buf := new(bytes.Buffer)

	if _, err := io.Copy(buf, resp.Body); err != nil {
		return http.StatusServiceUnavailable, ""
	}

	return http.StatusOK, buf.String()
}
