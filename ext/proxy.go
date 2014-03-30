package ext

import (
	"fmt"

	"github.com/j6n/noye/noye"
)

type ProxyBot struct{ bot noye.Bot }

func NewProxyBot(bot noye.Bot) *ProxyBot { return &ProxyBot{bot} }

func (p *ProxyBot) Reply(msg noye.Message, f string, a ...interface{}) {
	p.bot.Privmsg(msg.Target, fmt.Sprintf(msg.From+": "+f, a...))
}

func (p *ProxyBot) Send(f string, a ...interface{}) {
	p.bot.Send(f, a...)
}

func (p *ProxyBot) Privmsg(target, msg string) {
	p.bot.Privmsg(target, msg)
}

func (p *ProxyBot) Join(target string) {
	p.bot.Join(target)
}

func (p *ProxyBot) Part(target string) {
	p.bot.Part(target)
}
