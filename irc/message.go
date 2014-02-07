package irc

import "strings"

type Message struct {
	Source  string
	Command string
	Args    []string
	Text    string
}

func parse(raw string) Message {
	msg := Message{}
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
