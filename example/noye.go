package main

import (
	"log"

	"github.com/j6n/noye/irc"
)

func main() {
	log.Println("Starting connection")
	conn := &irc.Connection{}

	bot := irc.New(conn)
	if err := bot.Dial("localhost:6667", "noye"); err != nil {
		log.Fatalln(err)
	}

	<-bot.Wait()
	log.Println("done")
}
