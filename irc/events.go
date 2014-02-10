package irc

import (
	"strings"
	"sync"

	"github.com/j6n/noye/noye"
)

type Base struct {
	Bot noye.Bot
	cmd string
}

func NewBase(cmd string) *Base { return &Base{cmd: cmd} }

func (b *Base) Init(bot noye.Bot) { b.Bot = bot }
func (b *Base) Command() string   { return b.cmd }

type Ping struct{ *Base }

func PingEvent() *Ping { return &Ping{NewBase("PING")} }

func (p *Ping) Invoke(msg noye.IrcMessage) {
	p.Bot.Send("PONG %s", msg.Text)
}

type Ready struct {
	*Base
	autojoin []string
}

func ReadyEvent(chs ...string) *Ready { return &Ready{NewBase("001"), chs} }

func (r *Ready) Invoke(msg noye.IrcMessage) {
	for _, ch := range r.autojoin {
		r.Bot.Join(ch)
	}
}

type Privmsg struct {
	*Base
	plugins []noye.Plugin
	once    sync.Once // so we can init the bot reference once
}

func PrivmsgEvent(ps ...noye.Plugin) *Privmsg {
	return &Privmsg{Base: NewBase("PRIVMSG"), plugins: ps}
}

func (p *Privmsg) Invoke(msg noye.IrcMessage) {
	// set bot ref once it exists, IoC..
	p.once.Do(func() {
		for _, plugin := range p.plugins {
			plugin.Hook(p.Bot)
		}
	})

	out := noye.Message{
		From: strings.Split(msg.Source, "!")[0],
		Text: msg.Text,
	}

	switch msg.Args[0][0] {
	case '#', '&':
		out.Target = msg.Args[0]
	default:
		out.Target = out.From
	}

	go func(out noye.Message) {
		for _, plugin := range p.plugins {
			plugin.Listen() <- out
		}
	}(out)
}
