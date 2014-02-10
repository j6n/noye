package main

import (
	"log"

	"github.com/j6n/noye/irc"

	// plugins
	"github.com/j6n/noye/plugin/admin"
	"github.com/j6n/noye/plugin/naver"
	"github.com/j6n/noye/plugin/translate"
)

func main() {
	bot := irc.New(&irc.Connection{})
	bot.Autojoin = []string{"#museun", "#nanashin"}

	bot.AddPlugin(admin.New())
	bot.AddPlugin(naver.New())
	bot.AddPlugin(translate.New())

	if err := bot.Dial("irc.quakenet.org:6667", "noye", "museun"); err != nil {
		log.Fatalln(err)
	}

	<-bot.Wait()
	log.Println("done")
}
