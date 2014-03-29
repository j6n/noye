package irc

import (
	"fmt"
	"github.com/j6n/noye/noye"

	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

var tests = map[string]noye.IrcMessage{
	"local join": {
		Raw:     ":foo!bar@irc.localhost JOIN #foobar",
		Source:  "foo!bar@irc.localhost",
		Command: "JOIN",
		Args:    []string{"#foobar"},
	},
	"join": {
		Raw:     ":foo!bar@irc.localhost JOIN :#foobar",
		Source:  "foo!bar@irc.localhost",
		Command: "JOIN",
		Args:    []string{"#foobar"},
	},
	"privmsg": {
		Raw:     ":foo!bar@irc.localhost PRIVMSG #foobar :hello world",
		Source:  "foo!bar@irc.localhost",
		Command: "PRIVMSG",
		Args:    []string{"#foobar"},
		Text:    "hello world",
	},
	"part": {
		Raw:     ":foo!bar@irc.localhost PART #foobar :bye",
		Source:  "foo!bar@irc.localhost",
		Command: "PART",
		Args:    []string{"#foobar"},
		Text:    "bye",
	},
	"raw": {
		Raw:     "NOTICE AUTH :*** Checking Ident",
		Command: "NOTICE",
		Args:    []string{"AUTH"},
		Text:    "*** Checking Ident",
	},
	"ping": {
		Raw:     "PING :1594198849",
		Command: "PING",
		Args:    []string{},
		Text:    "1594198849",
	},
	"no text": {
		Raw:     ":irc.localhost 004 museun irc.localhost beware1.6.2 dgikoswx biklmnoprstv",
		Source:  "irc.localhost",
		Command: "004",
		Args:    []string{"museun", "irc.localhost", "beware1.6.2", "dgikoswx", "biklmnoprstv"},
		Text:    "",
	},
	"many colons": {
		Raw:     ":foo!bar@irc.localhost PRIVMSG hello :hello world :) for more :colons",
		Source:  "foo!bar@irc.localhost",
		Command: "PRIVMSG",
		Args:    []string{"hello"},
		Text:    "hello world :) for more :colons",
	},
}

func TestParser(t *testing.T) {
	Convey("Given a parser", t, func() {
		for k, v := range tests {
			Convey(fmt.Sprintf("It should parse '%s' message", k), func() {
				So(parse(v.Raw), ShouldResemble, v)
			})
		}
	})
}
