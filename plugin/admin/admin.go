package admin

import (
	"regexp"

	"github.com/j6n/noye/noye"
	"github.com/j6n/noye/plugin"
)

// Admin is a plugin which changes the bots state
type Admin struct {
	*plugin.Base
}

// New returns a new Admin Plugin
func New() *Admin {
	admin := &Admin{plugin.New("admin")}
	go admin.process()
	return admin
}

func (a *Admin) process() {
	chanMatcher := plugin.RegexMatcher(regexp.MustCompile("(#.*?)$"), true)
	whitelist := []string{"museun"} // TODO figure out how to get/update this

	opts := plugin.Options{
		Each:      true,
		Whitelist: whitelist,
	}

	join := plugin.Respond("join", opts, chanMatcher)
	part := plugin.Respond("part", opts, chanMatcher)

	handle := func(msg noye.Message) {
		switch {
		case join.Match(msg):
			for _, result := range join.Results() {
				a.Bot.Join(result)
			}
		case part.Match(msg):
			for _, result := range part.Results() {
				a.Bot.Part(result)
			}
		}
	}

	for msg := range a.Listen() {
		a.SafeHandle(handle, msg)
	}
}
