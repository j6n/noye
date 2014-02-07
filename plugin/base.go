package plugin

import "github.com/j6n/noye/noye"

type BasePlugin struct {
	Bot      noye.Bot
	Messages chan noye.Message
	Disabled map[string]bool
}

func New() *BasePlugin {
	return &BasePlugin{nil, make(chan noye.Message), make(map[string]bool)}
}

func (b *BasePlugin) Listen() chan noye.Message {
	return b.Messages
}

func (b *BasePlugin) Status(ch string) bool {
	if ch == "*" {
		for _, ok := range b.Disabled {
			if !ok {
				return false
			}
		}
	}

	s, ok := b.Disabled[ch]
	return s && ok
}

func (b *BasePlugin) SetStatus(ch string, ok bool) {
	if ch == "*" {
		for k, _ := range b.Disabled {
			b.Disabled[k] = ok
		}
		return
	}

	b.Disabled[ch] = ok
}

func (b *BasePlugin) Hook(bot noye.Bot) {
	b.Bot = bot
}
