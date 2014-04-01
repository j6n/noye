package ext

import (
	"fmt"
	"strings"

	"github.com/j6n/noye/noye"
)

// ProxyBot represents a noye.Bot with some parts of it no-opped
type ProxyBot struct{ bot noye.Bot }

// NewProxyBot returns a new bot, wrapping a noye.Bot
func NewProxyBot(bot noye.Bot) *ProxyBot { return &ProxyBot{bot} }

// Reply preprends the msg.From and sends it as a Privmsg to msg.Target
func (p *ProxyBot) Reply(msg noye.Message, f string, a ...interface{}) {
	out := fmt.Sprintf(msg.From+": "+f, a...)
	out = strings.Trim(out, "\r\n")
	p.bot.Privmsg(msg.Target, out)
}

// Send delegates to noye.Bot.Send
func (p *ProxyBot) Send(f string, a ...interface{}) {
	p.bot.Send(f, a...)
}

// Privmsg delegates to noye.Bot.Privmsg
func (p *ProxyBot) Privmsg(target, msg string) {
	p.bot.Privmsg(target, msg)
}

// Join delegates to noye.Bot.Join
func (p *ProxyBot) Join(target string) {
	p.bot.Join(target)
}

// Part delegates to noye.Bot.Part
func (p *ProxyBot) Part(target string) {
	p.bot.Part(target)
}
