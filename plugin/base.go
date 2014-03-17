package plugin

import (
	"fmt"
	"log"

	"github.com/j6n/noye/noye"
)

// Handler encapsulate a Command and gives it a callback
type Handler struct {
	*Command
	Handle func(*Command, noye.Message)
}

// Base is a type to reduce boilerplate of plugins
// It holds the noye.Bot reference, a channel to receive
// noye.Messages on and a map of disabled things
type Base struct {
	Bot   noye.Bot
	name  string
	input chan noye.Message

	Handlers []*Handler
}

// New returns a new Base plugin
func New(name string, handlers ...*Handler) *Base {
	base := &Base{
		name:     name,
		input:    make(chan noye.Message),
		Handlers: handlers,
	}

	go base.process()
	return base
}

// Listen returns the message to receives messages upon
func (b *Base) Listen() chan noye.Message {
	return b.input
}

// Name returns the plugins name
func (b *Base) Name() string {
	return b.name
}

// Hook sets the noye.Bot reference
func (b *Base) Hook(bot noye.Bot) {
	b.Bot = bot
}

// Reply replies to the noye.Message with a formatted string
func (b *Base) Reply(msg noye.Message, f string, a ...interface{}) {
	b.Bot.Privmsg(msg.Target, fmt.Sprintf(f, a...))
}

// Error replies to the noye.Message with a readable error
func (b *Base) Error(msg noye.Message, text string, err error) {
	b.Reply(msg, "error with %s (%s)", text, err)
}

func (b *Base) process() {
	for msg := range b.input {
		for _, cmd := range b.Handlers {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("recover! %s/%s from %s\n", b.Name(), cmd.Command.Command, err)
				}
			}()
			if cmd.Match(msg) {
				if !cmd.AcceptsFrom(msg.From) {
					b.Reply(msg, "You can't do this command.")
					continue
				}
				cmd.Handle(cmd.Command, msg)
			}
		}
	}
}
