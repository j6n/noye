package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func Follow(url string) (string, int) {
	client := &http.Client{}

	resp, err := client.Get(url)
	if err != nil {
		return "", http.StatusBadRequest
	}

	// 19+ won't redirect
	if resp.Request.URL.String() == url {
		return "", http.StatusUnauthorized
	}

	return resp.Request.URL.String(), http.StatusOK
}

func Get(url string, headers map[string]string) (string, int) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", http.StatusBadRequest
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", resp.StatusCode
	}

	defer resp.Body.Close()
	buf := new(bytes.Buffer)

	if _, err := io.Copy(buf, resp.Body); err != nil {
		return "", http.StatusServiceUnavailable
	}

	return buf.String(), http.StatusOK
}

const (
	googlURL = "https://www.googleapis.com/urlshortener/v1/url"
)

func Shorten(input string) (short string, status int) {
	client := &http.Client{}
	status = http.StatusOK

	// generate the json, make sure the input is sanitized
	reader := bytes.NewReader([]byte(fmt.Sprintf(`{"longUrl":"%s"}`, strings.Trim(input, " \r\n"))))

	// do a POST and get the response
	resp, err := client.Post(fmt.Sprintf(googlURL), "application/json", reader)
	if err != nil {
		status = http.StatusBadRequest
		return
	}

	defer resp.Body.Close()
	// read body into a []byte
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		status = http.StatusInternalServerError
		return
	}

	// just a simple generic map
	var j map[string]interface{}
	// deserialize []byte into the json
	if err = json.Unmarshal(data, &j); err != nil {
		status = http.StatusInternalServerError
		return
	}

	// if "code" is present, we've ran into an error
	if _, ok := j["code"].(string); ok {
		status = http.StatusInternalServerError
		return
	}

	// this is the 'shortened' url
	short = j["id"].(string)
	return
}
