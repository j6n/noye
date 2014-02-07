package admin

import (
	"log"

	"github.com/j6n/noye/dsl"
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
	join := dsl.Nick("noye").Command("join").Param("(#.*?)$")
	part := dsl.Nick("noye").Command("part").Param("(#.*?)$")

	// check to see if our join command is valid
	if ok, err := join.Valid(); !ok {
		log.Println("err starting admin:", err)
		return
	}

	// check to see if our part command is valid
	if ok, err := part.Valid(); !ok {
		log.Println("err starting admin:", err)
		return
	}

	// lazy, copy the ref to the bot
	ctx := a.BasePlugin.Bot

	// when we get a message
	for msg := range a.BasePlugin.Messages {
		switch {
		// see if its a join command
		case join.Match(msg):
			// if so join the channel
			ctx.Join(join.Results.Params()[0])

		// see if its a part command
		case part.Match(msg):
			// if so leave the channel
			ctx.Part(part.Results.Params()[0])
		}
	}
}
