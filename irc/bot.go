package irc

import (
	"fmt"
	"log"
	"sync"

	"github.com/j6n/noye/noye"
)

type Bot struct {
	conn noye.Conn
	stop chan struct{}

	events map[string][]noye.Event
	once   sync.Once
}

func New(conn noye.Conn) *Bot {
	bot := &Bot{
		conn:   conn,
		stop:   make(chan struct{}),
		events: make(map[string][]noye.Event),
	}

	return bot
}

func (b *Bot) Dial(addr, nick, user string) (err error) {
	if err = b.conn.Dial(addr); err != nil {
		return
	}

	b.Send("NICK %s", nick)
	b.Send("USER %s * 0 :%s", user, "noye in go!")

	go b.readLoop()
	return
}

func (b *Bot) Send(f string, a ...interface{}) {
	b.conn.WriteLine(fmt.Sprintf(f, a...))
}

func (b *Bot) Privmsg(t, msg string) { b.Send("PRIVMSG %s :%s", t, msg) }

func (b *Bot) Join(t string) { b.Send("JOIN %s", t) }
func (b *Bot) Part(t string) { b.Send("PART %s", t) }

func (b *Bot) Wait() <-chan struct{} { return b.stop }

func (b *Bot) AddEvent(ev noye.Event) {
	ev.Init(b)
	cmd := ev.Command()
	b.events[cmd] = append(b.events[cmd], ev)
}

func (b *Bot) Close() {
	b.once.Do(func() {
		b.conn.Close()
		close(b.stop)
	})
}

func (b *Bot) readLoop() {
	defer func() { b.Close() }()

	for {
		line, err := b.conn.ReadLine()
		if err != nil {
			log.Println(err)
			return
		}

		msg := parse(line)
		if evs, ok := b.events[msg.Command]; ok {
			for _, ev := range evs {
				ev.Invoke(msg)
			}
		}
	}
}
