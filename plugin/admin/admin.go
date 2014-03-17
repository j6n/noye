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

	join := &plugin.Handler{
		plugin.Respond("join", opts, chanMatcher),
		func(cmd *plugin.Command, msg noye.Message) {
			for _, result := range cmd.Results() {
				admin.Bot.Join(result)
			}
		},
	}

	part := &plugin.Handler{
		plugin.Respond("part", opts, chanMatcher),
		func(cmd *plugin.Command, msg noye.Message) {
			for _, result := range cmd.Results() {
				admin.Bot.Part(result)
			}
		},
	}

	admin.Base = plugin.New("admin", join, part)
	return admin
}
