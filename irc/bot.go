package irc

import (
	"fmt"
	"strings"
	"sync"

	"github.com/j6n/logger"
	"github.com/j6n/noye/noye"
)

var log = logger.New()

// Bot encapsulates all the parts to run a bot
type Bot struct {
	conn noye.Conn
	once sync.Once

	stop, ready *Signal
}

// New takes a noye.Conn and returns a new Bot
func New(conn noye.Conn) *Bot {
	bot := &Bot{
		conn:  conn,
		stop:  NewSignal(),
		ready: NewSignal(),
	}

	log.Debug("created a new bot")
	return bot
}

// Dial takes an address, nick and user string then connects and returns any error
func (b *Bot) Dial(addr, nick, user string) (err error) {
	fields := logger.Fields{"addr": addr, "nick": nick, "user": user}
	log.WithFields(fields).Info("connecting...")

	if err = b.conn.Dial(addr); err != nil {
		log.WithFields(fields).Errorf("while dialing: %v", err)
		return
	}

	b.Send("NICK %s", nick)
	b.Send("USER %s * 0 :%s", user, "noye in go!")

	log.Debug("starting read loop")
	go b.readLoop()
	return
}

// Send sends a formatted string to the connection
func (b *Bot) Send(f string, a ...interface{}) {
	msg := fmt.Sprintf(f, a...)
	log.WithField("send", true).Debug(msg)
	b.conn.WriteLine(msg)
}

// Privmsg sends the 'msg' to the target as a privmsg
func (b *Bot) Privmsg(t, msg string) {
	log.WithField("target", t).Info("sending message: %s", msg)
	b.Send("PRIVMSG %s :%s", t, msg)
}

// Join attempts to join the target
func (b *Bot) Join(t string) {
	log.WithField("target", t).Info("joining")
	b.Send("JOIN %s", t)
}

// Part attempts to leave the target
func (b *Bot) Part(t string) {
	log.WithField("target", t).Info("leaving")
	b.Send("PART %s", t)
}

// Quit closes the connection
func (b *Bot) Quit() {
	b.Send("QUIT %s", "bye")
}

// Wait returns a channel that'll be closed when the bot dies
func (b *Bot) Wait() <-chan struct{} {
	log.Debug("waiting to stop")

	return b.stop.Wait()
}

// Ready returns a channel that'll be closed when the bot is ready
func (b *Bot) Ready() <-chan struct{} {
	log.Debug("waiting to be ready")

	return b.ready.Wait()
}

// Close attempts to close the bots connection
func (b *Bot) Close() {
	log.Debug("attempting to close connection")
	b.once.Do(func() {
		log.Info("closing connection")
		b.conn.Close()
		b.stop.Close()
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
			log.Errorf("reading a line: %v", err)
			return
		}

		b.handle(line)
	}
}

func (b *Bot) handle(line string) {
	msg := parse(line)
	//	log.WithField("handle", msg.Command).Debug(msg.Raw)

	switch msg.Command {
	case "PING":
		b.Send("PONG %s", msg.Text)
	case "001":
		b.ready.Close()
		log.Info("connected")
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

func init() {
	log.Level = logger.Debug
	//log.Formatter = &logger.JsonFormatter{true}
}

func LoggerLevel(level logger.Level) {
	log.Level = level
}

func LoggerFormatter(format logger.Formatter) {
	log.Formatter = format
}
