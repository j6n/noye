package irc

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/j6n/logger"
	"github.com/j6n/noye/ext"
	"github.com/j6n/noye/noye"
)

var log = logger.New()
var Logger = log

// Bot encapsulates all the parts to run a bot
type Bot struct {
	conn noye.Conn
	once sync.Once

	manager *ext.Manager

	stop, ready *Signal
}

// New takes a noye.Conn and returns a new Bot
func New(conn noye.Conn) *Bot {
	bot := &Bot{
		conn:  conn,
		stop:  NewSignal(),
		ready: NewSignal(),
	}

	bot.manager = ext.New(bot, log)

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
	log.WithFields(logger.Fields{"target": t, "cmd": "privmsg", "data": msg}).Info("sending message")
	b.Send("PRIVMSG %s :%s", t, msg)
}

// Join attempts to join the target
func (b *Bot) Join(t string) {
	log.WithFields(logger.Fields{"target": t, "cmd": "join"}).Info("joining")
	b.Send("JOIN %s", t)
}

// Part attempts to leave the target
func (b *Bot) Part(t string) {
	log.WithFields(logger.Fields{"target": t, "cmd": "part"}).Info("leaving")
	b.Send("PART %s", t)
}

// Quit closes the connection
func (b *Bot) Quit() {
	log.WithField("cmd", "quit").Info("bye")
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

// Manager returns the script manager
func (b *Bot) Manager() noye.Manager {
	return b.manager
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

	// built-in switch
	switch msg.Command {
	case "PING":
		log.WithField("cmd", "ping").Debug(msg.Text)
		b.Send("PONG %s", msg.Text)
	case "001":
		log.Info("connected")
		b.ready.Close()
	case "PRIVMSG":
		out := ircToMsg(msg)
		b.manager.Respond(out)
	}

	// default listeners
	b.manager.Listen(msg)
}

func ircToMsg(msg noye.IrcMessage) noye.Message {
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

	fields := logger.Fields{
		"cmd":  "recv",
		"from": out.From,
	}

	if out.Target == msg.Args[0] {
		fields["type"] = "public"
		fields["target"] = out.Target
	} else {
		fields["type"] = "private"
	}
	return out
}

func init() {
	switch os.Getenv("NOYE_ENV") {
	case "TEST":
		log.Formatter = &logger.JsonFormatter{true}
	case "DEV":
		log.Formatter = &logger.TextFormatter{}
		log.Level = logger.Debug
	case "DIST":
		log.Formatter = &logger.JsonFormatter{false}
	}
}

// LoggerLevel sets the bots logger level
func LoggerLevel(level logger.Level) {
	log.Level = level
}

// LoggerFormatter sets the bots logger formatter
func LoggerFormatter(format logger.Formatter) {
	log.Formatter = format
}
