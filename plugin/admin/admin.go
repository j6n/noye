package admin

import (
	"regexp"

	"github.com/j6n/noye/noye"
	"github.com/j6n/noye/plugin"
)

// Admin is a plugin which changes the bots state
type Admin struct{ *plugin.Base }

// New returns a new Admin Plugin
func New() *Admin {
	admin := new(Admin)

	chanMatcher := plugin.RegexMatcher(regexp.MustCompile("(#.*?)$"), true)
	whitelist := []string{"museun"} // TODO figure out how to get/update this

	opts := plugin.Options{
		Each:      true,
		Whitelist: whitelist,
	}

	join := plugin.Respond("join", opts, chanMatcher)
	part := plugin.Respond("part", opts, chanMatcher)

	admin.Base = plugin.New("admin", join, part)
	admin.Handler = func(msg noye.Message) {
		switch {
		case join.Match(msg):
			for _, result := range join.Results() {
				admin.Bot.Join(result)
			}
		case part.Match(msg):
			for _, result := range part.Results() {
				admin.Bot.Part(result)
			}
		}
	}

	return admin
}
