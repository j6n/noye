package dsl

import (
	"regexp"
	"strings"

	"github.com/j6n/noye/noye"
)

type Matcher struct {
	Results Results

	nick *regexp.Regexp
	err  error

	cmds, params, lists []*regexp.Regexp
}

func New() *Matcher {
	return &Matcher{Results: Results{make(map[string][]string)}}
}

func Nick(nick string) *Matcher {
	matcher := &Matcher{Results: Results{make(map[string][]string)}}
	re, err := regexp.Compile(`(?:` + nick + `[:,]?\s*)`)
	if err != nil {
		matcher.err = err
		return nil
	}

	matcher.nick = re
	return matcher
}

func (m *Matcher) Command(cmd string) *Matcher {
	re, err := regexp.Compile(cmd)
	if err != nil {
		m.err = err
		return nil
	}

	m.cmds = append(m.cmds, re)
	return m
}

func (m *Matcher) Param(param string) *Matcher {
	re, err := regexp.Compile(param)
	if err != nil {
		m.err = err
		return nil
	}

	m.params = append(m.params, re)
	return m
}

func (m *Matcher) List(params ...string) *Matcher {
	for _, param := range params {
		re, err := regexp.Compile(param)
		if err != nil {
			m.err = err
			return nil
		}

		m.lists = append(m.lists, re)
	}

	return m
}

func (m *Matcher) Valid() (bool, error) {
	return m.err == nil, m.err
}

func (m *Matcher) Match(msg noye.Message) (ok bool) {
	// reset results
	m.Results = Results{make(map[string][]string)}

	params := strings.Fields(msg.Text)
	index := 0

	type match struct {
		re   *regexp.Regexp
		elem string
	}

	matches := make(map[string][]match)

	if m.nick != nil {
		if !m.nick.MatchString(params[index]) {
			return false
		}

		index++
	}

	// if we've gone too far
	if index >= len(params) {
		return false
	}

	for _, cmd := range m.cmds {
		if !cmd.MatchString(params[index]) {
			return false
		}

		matches["cmds"] = append(matches["cmds"], match{cmd, params[index]})
		index++
	}

	for _, param := range m.params {
		if !param.MatchString(params[index]) {
			return false
		}

		matches["params"] = append(matches["params"], match{param, params[index]})
		index++
	}

	for _, param := range params[index:] {
		for _, list := range m.lists {
			if !list.MatchString(param) {
				return false
			}

			matches["lists"] = append(matches["lists"], match{list, param})
		}
	}

	for k, v := range matches {
		for _, p := range v {
			match := p.re.FindStringSubmatch(p.elem)[1:]
			m.Results.data[k] = append(m.Results.data[k], match...)
		}
	}

	return true
}

type Results struct{ data map[string][]string }

func (r *Results) Cmds() []string   { return r.data["cmds"] }
func (r *Results) Params() []string { return r.data["params"] }
func (r *Results) Lists() []string  { return r.data["lists"] }
