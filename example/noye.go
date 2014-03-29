package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"

	"github.com/j6n/noye/irc"
)

func init() {
	runtime.GOMAXPROCS(4)
}

func main() {
	// to capture Ctrl-C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// set logger to info
	//irc.LoggerLevel(logger.Info)

	bot := irc.New(&irc.Connection{})

	if err := bot.Dial("irc.quakenet.org:6667", "noye", "museun"); err != nil {
		log.Fatalln(err)
	}

	go func() {
		<-quit
		// send the quit
		bot.Quit()
	}()

	// wait for the ready signal
	<-bot.Ready()
	bot.Join("#noye")

	// wait for the close signal
	<-bot.Wait()
	log.Println("done")
}
