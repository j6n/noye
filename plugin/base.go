package plugin

import (
	"fmt"

	"github.com/j6n/noye/noye"
)

type BasePlugin struct {
	Bot      noye.Bot
	input    chan noye.Message
	Disabled map[string]bool
}

func New() *BasePlugin {
	return &BasePlugin{
		input:    make(chan noye.Message),
		Disabled: make(map[string]bool),
	}
}

func (b *BasePlugin) Listen() chan noye.Message { return b.input }

func (b *BasePlugin) Status(ch string) bool {
	s, ok := b.Disabled[ch]
	return s && ok
}

func (b *BasePlugin) SetStatus(ch string, status bool) { b.Disabled[ch] = status }

func (b *BasePlugin) Hook(bot noye.Bot) { b.Bot = bot }

func (b *BasePlugin) Reply(msg noye.Message, f string, a ...interface{}) {
	b.Bot.Privmsg(msg.Target, fmt.Sprintf(f, a...))
}

func (b *BasePlugin) Error(msg noye.Message, text string, err error) {
	b.Reply(msg, "error with %s (%s)", text, err)
}
