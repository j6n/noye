package plugin

import "github.com/j6n/noye/noye"

type Plugin interface {
	// return a channel to receive PMs on
	Listen() chan noye.Message
	// gets status for channel, * for all
	Status(string) bool
	// sets status for channel, * for all
	SetStatus(string, bool)
	// sets the plugin context
	Hook(noye.Bot)
}
