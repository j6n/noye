package plugin

import (
	"fmt"
	"log"

	"github.com/j6n/noye/noye"
)

// Base is a type to reduce boilerplate of plugins
// It holds the noye.Bot reference, a channel to receive
// noye.Messages on and a map of disabled things
type Base struct {
	Bot   noye.Bot
	name  string
	input chan noye.Message

	Disabled map[string]bool
	Commands map[string]*Command
	Handler  func(noye.Message)
}

// New returns a new Base plugin
func New(name string, cmds ...*Command) *Base {
	base := &Base{
		name:  name,
		input: make(chan noye.Message),

		Disabled: make(map[string]bool),
		Commands: make(map[string]*Command),
		Handler:  func(noye.Message) {},
	}

	for _, cmd := range cmds {
		base.Commands[cmd.Command] = cmd
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

// Status returns whether the input is disabled
func (b *Base) Status(ch string) bool {
	s, ok := b.Disabled[ch]
	return s && ok
}

// SetStatus disables the input
func (b *Base) SetStatus(ch string, status bool) {
	b.Disabled[ch] = status
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
		defer func() {
			if err := recover(); err != nil {
				log.Println("recover!", b.Name(), "from", err)
			}
		}()

		b.Handler(msg)
	}
}
