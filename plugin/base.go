package plugin

import (
	"fmt"

	"github.com/j6n/noye/noye"
)

type BasePlugin struct {
	Bot      noye.Bot
	Messages chan noye.Message
	Disabled map[string]bool
}

func New() *BasePlugin {
	return &BasePlugin{
		Messages: make(chan noye.Message),
		Disabled: make(map[string]bool),
	}
}

func (b *BasePlugin) Listen() chan noye.Message {
	return b.Messages
}

func (b *BasePlugin) Status(ch string) (ok bool) {
	if ch == "*" {
		for _, ok := range b.Disabled {
			if !ok {
				return
			}
		}
	}

	s, ok := b.Disabled[ch]
	return s && ok
}

func (b *BasePlugin) SetStatus(ch string, status bool) {
	if ch == "*" {
		for k, _ := range b.Disabled {
			b.Disabled[k] = status
		}
		return
	}

	b.Disabled[ch] = status
}

func (b *BasePlugin) Hook(bot noye.Bot) {
	b.Bot = bot
}

func (b *BasePlugin) Reply(msg noye.Message, f string, a ...interface{}) {
	b.Bot.Privmsg(msg.Target, fmt.Sprintf(f, a...))
}

func (b *BasePlugin) Error(msg noye.Message, text string, err error) {
	b.Reply(msg, "error with %s (%s)", text, err)
}
