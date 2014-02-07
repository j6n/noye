package irc

import (
	"fmt"
	"log"
	"sync"
)

type Bot struct {
	conn Conn
	stop chan struct{}
	once sync.Once
}

func New(conn Conn) *Bot {
	bot := &Bot{
		conn: conn,
		stop: make(chan struct{}),
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
	log.Println(">", msg)
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
			b.Send("JOIN %s", "#test")
		case "PRIVMSG":
			b.handle(msg)
		}
	}
}

func (b *Bot) handle(msg IrcMessage) {
	log.Printf("%q\n", msg)
}
