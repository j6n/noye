package script

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/robertkrimen/otto"
)

func safeRun(fn scriptFunc, name string, vals ...otto.Value) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()

	fn(vals[0], vals[1:]...)
}

func findMatches(s string, re *regexp.Regexp) map[string]string {
	captures := make(map[string]string)
	match := re.FindStringSubmatch(s)
	if match == nil {
		return captures
	}

	for i, name := range re.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}

		captures[name] = strings.TrimSpace(match[i])
	}

	return captures
}

func valuesToInterfaces(vals []otto.Value) (out []interface{}) {
	for _, val := range vals {
		if res, err := val.Export(); err == nil {
			out = append(out, res)
		}
	}
	return
}

func valuesToString(vals []otto.Value) string {
	var msg []string
	for _, val := range vals {
		if it, err := val.Export(); err == nil {
			msg = append(msg, fmt.Sprintf("%+v", it))
		}
	}

	return strings.Join(msg, " ")
}

func checkAnyError(errs ...error) (err error) {
	for _, err = range errs {
		if err != nil {
			return
		}
	}

	return
}
