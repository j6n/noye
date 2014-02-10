package urlinfo

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/j6n/noye/noye"
	"github.com/j6n/noye/plugin"
)

type Urlinfo struct {
	*plugin.BasePlugin
	blacklist map[string]*regexp.Regexp
}

func New() *Urlinfo {
	urlinfo := &Urlinfo{
		plugin.New(),
		make(map[string]*regexp.Regexp),
	}
	go urlinfo.process()
	return urlinfo
}

var urlinfoRe = regexp.MustCompile(`^(\w+:\/\/[\w@][\w.:@]+\/?[\w\.?=%&=\-@/$,:]*)$`)

func (u *Urlinfo) process() {
	urlinfo := plugin.Command{Matcher: plugin.RegexMatcher(urlinfoRe, true), Each: true}
	add := plugin.Respond("url-block", plugin.NoopMatcher())
	del := plugin.Respond("url-unblock", plugin.NoopMatcher())

	for msg := range u.Listen() {
		switch {
		case urlinfo.Match(msg):
			u.handleUrl(msg, urlinfo.Results())
		case add.Match(msg):
			u.handleAdd(msg, urlinfo.Results())
		case del.Match(msg):
			u.handleDel(msg, urlinfo.Results())
		}
	}
}

func (u *Urlinfo) handleAdd(msg noye.Message, results []string) {

}

func (u *Urlinfo) handleDel(msg noye.Message, results []string) {

}

func (u *Urlinfo) handleUrl(msg noye.Message, results []string) {
	for _, part := range results {
		url := strings.TrimSpace(urlinfoRe.FindString(part))
		resp, err := client.Get(url)
		if err != nil {
			// can't get page
			if n, ok := err.(net.Error); ok && n.Timeout() {
				// timedout
			}
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode > 300 {
			// not a status-ok result
			continue
		}

		body := make([]byte, 10*1024)
		if _, err := io.ReadAtLeast(resp.Body, body, 10*1024); err != nil {
			continue
		}

		if title, ok := findTitle(body); ok {
			log.Printf("%s\n", title)
			continue
		}
	}
}

func findTitle(body []byte) (title []byte, ok bool) {
	switch start, end := bytes.Index(body, []byte("<title>")), bytes.Index(body, []byte("</title>")); {
	case start > -1 && end > -1:
		return body[start+bytes.Index(body[start:end], []byte(">"))+1 : end], true
	case start > -1:
		return body[start+bytes.Index(body[start:], []byte(">"))+1:], true
	default:
		return
	}
}

var (
	client *http.Client
)

func init() {
	client = &http.Client{Transport: &http.Transport{Dial: timeout(2*time.Second, 2*time.Second)}}
}

func timeout(min, max time.Duration) func(string, string) (net.Conn, error) {
	return func(t, addr string) (conn net.Conn, err error) {
		if conn, err = net.DialTimeout(t, addr, min); err != nil {
			return
		}

		conn.SetDeadline(time.Now().Add(max))
		return
	}
}
