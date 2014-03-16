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
	Bot      noye.Bot
	name     string
	input    chan noye.Message
	Disabled map[string]bool
}

// New returns a new Base plugin
func New(name string) *Base {
	return &Base{
		name:     name,
		input:    make(chan noye.Message),
		Disabled: make(map[string]bool),
	}
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

// SafeHandler is a function that'll safely execute and capture any panics
type SafeHandler func(noye.Message)

// SafeHandle takes a SafeHandler and a msg and captures any panics
func (b *Base) SafeHandle(fn SafeHandler, msg noye.Message) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("recover!", b.Name(), "from", err)
		}
	}()

	fn(msg)
}
