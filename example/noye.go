package main

import (
	"log"

	"github.com/j6n/noye/irc"
)

func main() {
	bot := irc.New(&irc.Connection{})
	if err := bot.Dial("irc.quakenet.org:6667", "noye", "museun"); err != nil {
		log.Fatalln(err)
	}

	<-bot.Wait()
	log.Println("done")
}
