package main

import (
	"log"

	"github.com/j6n/noye/irc"
	"github.com/j6n/noye/plugin/admin"
	"github.com/j6n/noye/plugin/naver"
)

func main() {
	bot := irc.New(&irc.Connection{})
	bot.Autojoin = []string{"#museun", "#nanashin"}

	bot.AddPlugin(admin.New())
	bot.AddPlugin(naver.New())

	if err := bot.Dial("irc.quakenet.org:6667", "noye", "museun"); err != nil {
		log.Fatalln(err)
	}

	<-bot.Wait()
	log.Println("done")
}
