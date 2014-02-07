package irc

import (
	"fmt"
	"log"
	"sync"
)

type Bot struct {
	Handle   func(msg IrcMessage)
	Autojoin []string

	conn Conn
	stop chan struct{}
	once sync.Once
}

func New(conn Conn) *Bot {
	bot := &Bot{
		conn: conn,
		stop: make(chan struct{}),

		Autojoin: make([]string, 0),
		Handle:   func(msg IrcMessage) {},
	}

	return bot
}

func (b *Bot) Dial(addr, nick string) (err error) {
	if err = b.conn.Dial(addr, nick); err != nil {
		return
	}

	b.Send("NICK %s", nick)
	b.Send("USER %s * 0 :%s", nick, nick)

	go b.readLoop()
	return
}

func (b *Bot) Send(f string, a ...interface{}) {
	msg := fmt.Sprintf(f, a...)
	b.conn.WriteLine(msg)
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
			b.Handle(msg)
		}
	}
}
