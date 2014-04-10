package irc

import (
	"strings"

	"github.com/j6n/noye/noye"
)

// Parse takes a raw message and turns it into an IrcMessage
func Parse(raw string) noye.IrcMessage {
	msg := noye.IrcMessage{Raw: raw}

	// :source command [args] :message
	if raw[0] == ':' {
		if i := strings.Index(raw, " "); i >= -1 {
			msg.Source = ParseUser(raw[1:i])
			raw = raw[i+1 : len(raw)]
		}
	}

	args := strings.SplitN(raw, " :", 2)
	if len(args) > 1 {
		msg.Text = args[1]
	}

	args = strings.Split(args[0], " ")
	msg.Command, msg.Args = strings.ToUpper(args[0]), []string{}

	if len(args) > 1 {
		msg.Args = args[1:len(args)]
	}

	return msg
}

// ParseUser takes a raw user string and returns a User
func ParseUser(raw string) noye.User {
	if strings.Index(raw, "!") == -1 {
		return noye.User{Nick: raw}
	}

	first := strings.Split(raw, "!")
	nick, second := first[0], strings.Split(first[1], "@")
	user, host := second[0], second[1]

	return noye.User{nick, user, host}
}

func ircToMsg(msg noye.IrcMessage) noye.Message {
	out := noye.Message{From: msg.Source, Text: msg.Text}

	switch msg.Args[0][0] {
	case '#', '&':
		out.Target = msg.Args[0]
	default:
		out.Target = out.From.Nick
	}

	return out
}
