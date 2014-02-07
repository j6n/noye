package irc

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/j6n/noye/noye"
)

type Bot struct {
	Autojoin []string
	plugins  []noye.Plugin

	conn noye.Conn
	stop chan struct{}
	once sync.Once
}

func New(conn noye.Conn) *Bot {
	bot := &Bot{
		conn: conn,
		stop: make(chan struct{}),

		Autojoin: make([]string, 0),
	}

	return bot
}

func (b *Bot) Dial(addr, nick, user string) (err error) {
	if err = b.conn.Dial(addr, nick); err != nil {
		return
	}

	b.Send("NICK %s", nick)
	b.Send("USER %s * 0 :%s", user, "noye in go!")

	go b.readLoop()
	return
}

func (b *Bot) Send(f string, a ...interface{}) {
	msg := fmt.Sprintf(f, a...)
	b.conn.WriteLine(msg)
}

func (b *Bot) Privmsg(target, msg string) {
	b.Send("PRIVMSG %s :%s", target, msg)
}

func (b *Bot) Join(target string) {
	b.Send("JOIN %s", target)
}

func (b *Bot) Part(target string) {
	b.Send("PART %s", target)
}

func (b *Bot) Close() {
	b.once.Do(func() {
		b.conn.Close()
		close(b.stop)
	})
}

func (b *Bot) Wait() <-chan struct{} {
	return b.stop
}

func (b *Bot) AddPlugin(p noye.Plugin) {
	p.Hook(b)
	b.plugins = append(b.plugins, p)
}

func (b *Bot) readLoop() {
	defer func() { b.Close() }()

	for {
		line, err := b.conn.ReadLine()
		if err != nil {
			log.Println(err)
			return
		}

		switch msg := parse(line); msg.Command {
		case "PING":
			b.conn.WriteLine("PONG " + msg.Text)
		case "001":
			for _, join := range b.Autojoin {
				b.Send("JOIN %s", join)
			}
		case "PRIVMSG":
			b.handle(msg)
		}
	}
}

func (b *Bot) handle(msg Message) {
	out := noye.Message{}

	out.From = strings.Split(msg.Source, "!")[0]
	out.Text = msg.Text

	switch msg.Args[0][0] {
	case '#', '&':
		out.Target = msg.Args[0]
	default:
		out.Target = out.From
	}

	go func(out noye.Message) {
		for _, p := range b.plugins {
			p.Listen() <- out
		}
	}(out)
}
