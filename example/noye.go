package main

import (
	"log"

	"github.com/j6n/noye/irc"
	"github.com/j6n/noye/plugin/admin"
)

func main() {
	bot := irc.New(&irc.Connection{})
	bot.Autojoin = []string{"#test"}
	bot.AddPlugin(admin.New())

	if err := bot.Dial("localhost:6667", "noye", "museun"); err != nil {
		log.Fatalln(err)
	}

	<-bot.Wait()
	log.Println("done")
}
