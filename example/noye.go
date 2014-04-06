package main

import (
	"os"
	"os/signal"
	"runtime"

	"github.com/j6n/noye/config"
	"github.com/j6n/noye/store"

	"github.com/j6n/noye/irc"
	"github.com/j6n/noye/logger"
)

var (
	log  = logger.Get()
	conf = config.NewConfig()
	db   *store.DB
)

func init() {
	runtime.GOMAXPROCS(4)

	db = store.NewDB()
	if err := db.CheckTable("config", store.KvSchema); err != nil {
		log.Fatalf("can't create table %s:%s\n", "config", err)
	}

	m := conf.ToMap()
	for k, v := range m {
		db.Set("config", k, v)
	}
}

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	bot := irc.New(&irc.Connection{})

	reconnect := true
	go func() { <-quit; reconnect = false; bot.Quit() }()

	for reconnect {
		if err := bot.Dial(conf.Server, conf.Nick, conf.Nick); err != nil {
			return
		}

		<-bot.Wait()
	}
}
