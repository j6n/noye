package main

import (
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/j6n/noye/ext"

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

	for k, v := range conf.ToMap() {
		db.Set("config", k, v)
	}
}

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	store.Debug = true

	bot := irc.New(&irc.Connection{})
	bot.Manager().LoadScripts("./scripts")

	reconnect := true
	go func() { <-quit; reconnect = false; bot.Quit() }()

	defer func() {
		for _, script := range bot.Manager().Scripts() {
			script.Cleanup()
		}

		db.Close()
		<-time.After(3 * time.Second)
	}()

	for reconnect {
		if err := bot.Dial(conf.Server, conf.Nick, conf.User); err != nil {
			return
		}

		<-bot.Ready()
		broadcastConfig()

		<-bot.Wait()
	}
}

func broadcastConfig() {
	for key, val := range conf.ToMap() {
		ext.Broadcast(key, val)
	}
}
