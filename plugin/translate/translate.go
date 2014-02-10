package translate

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/j6n/noye/noye"
	"github.com/j6n/noye/plugin"
)

type Translate struct {
	*plugin.BasePlugin
}

func New() *Translate {
	translate := &Translate{plugin.New()}
	go translate.process()
	return translate
}

var translateRe = regexp.MustCompile(`(?P<to>[a-zA-Z-]+),(?P<from>[a-zA-Z-]+) (?P<text>.+?)$`)

func (t *Translate) process() {
	translate := plugin.Respond("translate", plugin.RegexMatcher(translateRe, false))

	for msg := range t.Listen() {
		switch {
		case translate.Match(msg):
			t.handleTranslate(msg)
		}
	}
}

func (t *Translate) handleTranslate(msg noye.Message) {
	matches := translateRe.FindAllStringSubmatch(msg.Text, -1)[0][1:]
	if len(matches) != 3 {
		t.Error(msg, "You didn't provide the right params: translate from,to message", nil)
		return
	}

	// TODO auto
	from, to, text := matches[0], matches[1], matches[2]

	resp, err := http.Get(baseUrl + url.Values{
		"client": {"a"}, "ie": {"UTF-8"}, "oe": {"UTF-8"}, "sc": {"1"},
		"q": {text}, "sl": {from}, "tl": {to}, "uptl": {to},
	}.Encode())
	if err != nil {
		t.Error(msg, "I can't do that", err)
		return
	}

	defer resp.Body.Close()
	buf := &bytes.Buffer{}
	io.Copy(buf, resp.Body)

	var res struct{ Sentences []struct{ Trans string } }
	if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
		t.Error(msg, "I can't do that", err)
		return
	}

	s := res.Sentences
	if len(s) == 0 {
		t.Error(msg, "I can't do that", nil)
		return
	}

	var b bytes.Buffer
	for _, v := range s {
		b.WriteString(v.Trans)
	}

	t.Reply(msg, b.String())
}

const (
	baseUrl = "https://translate.google.com/translate_a/t?"
)

var (
	LANG = [52]string{
		"af", "sq", "ar", "be", "bg", "ca", "zh-CN", "zh-TW", "hr",
		"cs", "da", "nl", "en", "et", "tl", "fi", "fr", "gl", "de",
		"el", "iw", "hi", "hu", "is", "id", "ga", "it", "ja", "ko",
		"lv", "lt", "mk", "ms", "mt", "no", "fa", "pl", "pt", "ro",
		"ru", "sr", "sk", "sl", "es", "sw", "sv", "th", "tr", "uk",
		"vi", "cy", "yi",
	}
)
