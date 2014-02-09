package admin

import (
	"regexp"

	"github.com/j6n/noye/plugin"
)

type Admin struct {
	*plugin.BasePlugin
}

func New() *Admin {
	admin := &Admin{plugin.New()}
	// start the main loop
	go admin.process()
	return admin
}

func (a *Admin) process() {
	// create a matcher for our commands
	chanMatcher := plugin.RegexMatcher(
		regexp.MustCompile("(#.*?)$"),
		true,
	)

	// create our commands
	join := plugin.Respond("join", chanMatcher)
	join.Each = true
	part := plugin.Respond("part", chanMatcher)
	part.Each = true

	// when we get a message
	for msg := range a.Listen() {
		switch {
		// see if its a join command
		case join.Match(msg):
			// if so join the channels
			for _, result := range join.Results() {
				a.Bot.Join(result)
			}

		// see if its a part command
		case part.Match(msg):
			// if so leave the channel
			for _, result := range part.Results() {
				a.Bot.Part(result)
			}
		}
	}
}
