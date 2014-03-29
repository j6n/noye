package irc

import (
	"strings"

	"github.com/j6n/noye/noye"
)

// parse takes a raw string and returns an IrcMessage
// by parsing it somewhat accordingly to the IRC RFC
func parse(raw string) noye.IrcMessage {
	msg := noye.IrcMessage{Raw: raw}

	// :source command [args] :message
	if raw[0] == ':' {
		if i := strings.Index(raw, " "); i >= -1 {
			msg.Source = raw[1:i]
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
