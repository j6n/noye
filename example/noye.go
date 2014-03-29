package main

import (
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

	//irc.LoggerFormatter(&logger.JsonFormatter{true})

	bot := irc.New(&irc.Connection{})
	if err := bot.Dial("irc.quakenet.org:6667", "noye", "museun"); err != nil {
		return
	}

	go func() {
		// wait for the ctrl-c
		<-quit
		// send the quit
		bot.Quit()
	}()

	// wait for the ready signal
	<-bot.Ready()
	bot.Join("#noye")

	// wait for the close signal
	<-bot.Wait()
}
