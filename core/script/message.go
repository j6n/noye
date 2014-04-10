package script

import (
	"fmt"
	"strings"

	"github.com/j6n/noye/noye"
)

type wrappedMessage struct {
	noye.Message
	context noye.Bot

	Public bool
}

func (m *Manager) wrapMessage(msg noye.Message) wrappedMessage {
	return wrappedMessage{
		Message: msg,
		context: m.context,
		Public:  msg.From.Nick != msg.Target,
	}
}

// Reply to the messages, with the senders name prepended
func (w wrappedMessage) Reply(f string, a ...interface{}) {
	w.Send("%s: %s", w.From.Nick, fmt.Sprintf(f, a...))
}

// Send to the target of the message
func (w wrappedMessage) Send(f string, a ...interface{}) {
	w.context.Privmsg(w.Target, strings.Trim(fmt.Sprintf(f, a...), "\r\n"))
}
