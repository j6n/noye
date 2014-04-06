package irc

import (
	"fmt"

	"github.com/j6n/noye/ext"
	"github.com/j6n/noye/logger"
	"github.com/j6n/noye/noye"
)

var log = logger.Get()

// Bot encapsulates all the parts to run a bot
type Bot struct {
	conn    noye.Conn
	manager *ext.Manager

	stop, ready *Signal
}

// New takes a noye.Conn and returns a new Bot
func New(conn noye.Conn) *Bot {
	bot := &Bot{conn: conn}
	bot.manager = ext.New(bot)
	return bot
}

// Dial takes an address, nick and user string then connects and returns any error
func (b *Bot) Dial(addr, nick, user string) (err error) {
	log.Infof("Connecting to '%s' with '%s,%s'\n", addr, nick, user)
	if err = b.conn.Dial(addr); err != nil {
		log.Errorf("Failed to connect to '%s': %s\n", err)
		return
	}

	b.stop = NewSignal()
	b.ready = NewSignal()

	b.Send("NICK %s", nick)
	b.Send("USER %s * 0 :%s", user, "noye in go!")

	go b.readLoop()
	return
}

// Send sends a formatted string to the connection
func (b *Bot) Send(f string, a ...interface{}) {
	msg := fmt.Sprintf(f, a...)
	log.Debugf(">> %s\n", msg)
	b.conn.WriteLine(msg)
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

// Quit closes the connection
func (b *Bot) Quit() {
	b.Send("QUIT %s", "bye")
}

// Wait returns a channel that'll be closed when the bot dies
func (b *Bot) Wait() <-chan struct{} {
	return b.stop.Wait()
}

// Ready returns a channel that'll be closed when the bot is ready
func (b *Bot) Ready() <-chan struct{} {
	return b.ready.Wait()
}

// Close attempts to close the bots connection
func (b *Bot) Close() {
	log.Debugf("Closing the bot\n")
	b.conn.Close()
}

// Manager returns the script manager
func (b *Bot) Manager() noye.Manager {
	return b.manager
}

func (b *Bot) readLoop() {
	defer func() { b.Close(); b.stop.Close() }()

	var (
		line string
		err  error
	)

	for {
		line, err = b.conn.ReadLine()
		if err != nil {
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
		b.Send("PONG %s", msg.Text)
	case "001":
		b.ready.Close()
	case "PRIVMSG":
		b.manager.Respond(ircToMsg(msg))
	}

	// default listeners
	b.manager.Listen(msg)
}
