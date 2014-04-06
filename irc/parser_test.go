package irc

import (
	"fmt"

	"github.com/j6n/noye/noye"

	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestParser(t *testing.T) {
	user := parseUser("foo!bar@irc.localhost")
	tests := map[string]noye.IrcMessage{
		"local join": {
			Raw:     ":foo!bar@irc.localhost JOIN #foobar",
			Source:  user,
			Command: "JOIN",
			Args:    []string{"#foobar"},
		},
		"join": {
			Raw:     ":foo!bar@irc.localhost JOIN :#foobar",
			Source:  user,
			Command: "JOIN",
			Args:    []string{"#foobar"},
		},
		"privmsg": {
			Raw:     ":foo!bar@irc.localhost PRIVMSG #foobar :hello world",
			Source:  user,
			Command: "PRIVMSG",
			Args:    []string{"#foobar"},
			Text:    "hello world",
		},
		"part": {
			Raw:     ":foo!bar@irc.localhost PART #foobar :bye",
			Source:  user,
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
			Source:  parseUser("irc.localhost"),
			Command: "004",
			Args:    []string{"museun", "irc.localhost", "beware1.6.2", "dgikoswx", "biklmnoprstv"},
			Text:    "",
		},
		"many colons": {
			Raw:     ":foo!bar@irc.localhost PRIVMSG hello :hello world :) for more :colons",
			Source:  user,
			Command: "PRIVMSG",
			Args:    []string{"hello"},
			Text:    "hello world :) for more :colons",
		},
	}

	Convey("Given a parser", t, func() {
		for k, v := range tests {
			Convey(fmt.Sprintf("It should parse '%s' message", k), func() {
				So(parse(v.Raw), ShouldResemble, v)
			})
		}
	})
}

func TestConvert(t *testing.T) {
	type message struct {
		in  noye.IrcMessage
		out noye.Message
	}
	user := parseUser("foo!bar@irc.localhost")
	tests := map[string]message{
		"private": {
			in: parse(":foo!bar@irc.localhost PRIVMSG hello :hello world"),
			out: noye.Message{
				From:   user,
				Target: "foo",
				Text:   "hello world",
			},
		},
		"public": {
			in: parse(":foo!bar@irc.localhost PRIVMSG #hello :hello world"),
			out: noye.Message{
				From:   user,
				Target: "#hello",
				Text:   "hello world",
			},
		},
	}

	Convey("For a set of messages", t, func() {
		for k, msg := range tests {
			Convey(fmt.Sprintf("It should convert '%s'", k), func() {
				So(ircToMsg(msg.in), ShouldResemble, msg.out)
			})
		}
	})
}
