package irc

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/j6n/noye/noye"
)

// Bot encapsulates all the parts to run a bot
type Bot struct {
	conn noye.Conn
	once sync.Once

	stop, ready chan struct{}
}

// New takes a noye.Conn and returns a new Bot
func New(conn noye.Conn) *Bot {
	bot := &Bot{
		conn:  conn,
		stop:  make(chan struct{}),
		ready: make(chan struct{}),
	}

	return bot
}

// Dial takes an address, nick and user string then connects and returns any error
func (b *Bot) Dial(addr, nick, user string) (err error) {
	if err = b.conn.Dial(addr); err != nil {
		return
	}

	b.Send("NICK %s", nick)
	b.Send("USER %s * 0 :%s", user, "noye in go!")

	go b.readLoop()
	return
}

// Send sends a formatted string to the connection
func (b *Bot) Send(f string, a ...interface{}) {
	b.conn.WriteLine(fmt.Sprintf(f, a...))
}

// Privmsg sends the 'msg' to the target as a privmsg
func (b *Bot) Privmsg(t, msg string) {
	b.Send("PRIVMSG %s :%s", t, msg)
}

// Join attempts to join the target
func (b *Bot) Join(t string) {
	b.Send("JOIN %s", t)
}

// Part attempts to leave the target
func (b *Bot) Part(t string) {
	b.Send("PART %s", t)
}

// Wait returns a channel that'll be closed when the bot dies
func (b *Bot) Wait() <-chan struct{} {
	return b.stop
}

// Ready returns a channel that'll be closed when the bot is ready
func (b *Bot) Ready() <-chan struct{} {
	return b.ready
}

// Close attempts to close the bots connection
func (b *Bot) Close() {
	b.once.Do(func() {
		b.conn.Close()
		close(b.stop)
	})
}

func (b *Bot) readLoop() {
	defer func() { b.Close() }()

	var (
		line string
		err  error
	)

	for {
		line, err = b.conn.ReadLine()
		if err != nil {
			log.Println(err)
			return
		}

		b.handle(line)
	}
}

func (b *Bot) handle(line string) {
	msg := parse(line)
	switch msg.Command {
	case "PING":
		// send a PONG back
		b.Send("PONG %s", msg.Text)
	case "001":
		// BUG: this should be using a sync.Once
		close(b.ready)
	case "PRIVMSG":
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
	}

	// default should delegate to any extra eventsn
}
