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
	// create our commands
	join := plugin.Command{
		Respond: true,
		Each:    true,
		Command: "join",
		Matcher: plugin.RegexMatcher(
			regexp.MustCompile("($.*?)$"),
			true,
		),
	}

	part := plugin.Command{
		Respond: true,
		Each:    true,
		Command: "part",
		Matcher: plugin.RegexMatcher(
			regexp.MustCompile("($.*?)$"),
			true,
		),
	}

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
