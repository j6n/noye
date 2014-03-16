package main

import (
	"log"
	"runtime"

	"github.com/j6n/noye/irc"

	// plugins
	"github.com/j6n/noye/plugin/admin"
)

func init() {
	runtime.GOMAXPROCS(4)
}

func main() {
	bot := irc.New(&irc.Connection{})

	bot.AddPlugin(admin.New())

	if err := bot.Dial("irc.quakenet.org:6667", "noye", "museun"); err != nil {
		log.Fatalln(err)
	}

	<-bot.Ready()
	bot.Join("#museun")

	<-bot.Wait()
	log.Println("done")
}
