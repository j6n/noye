package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/j6n/noye/core/config"
	"github.com/j6n/noye/core/irc"
	"github.com/j6n/noye/core/logger"
	"github.com/j6n/noye/core/store"
)

var (
	log  = logger.Get()
	conf = config.NewConfig()

	db    *store.DB
	debug bool
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

	var qpass, quser string
	if qpass, quser = os.Getenv("NOYE_PASS"), os.Getenv("NOYE_USER"); qpass != "" && quser != "" {
		m := map[string]string{"user": quser, "pass": qpass}
		data, _ := json.Marshal(m)
		db.Set("config", "quakenet", string(data))
	}

	if os.Getenv("NOYE_DEBUG") != "" {
		log.SetMinMaxSeverity(logger.TRACE, logger.PANIC)
		debug = true
	} else {
		log.SetMinMaxSeverity(logger.INFO, logger.PANIC)
	}
}

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	if debug {
		store.Debug = true
	}

	bot := irc.New(&irc.Connection{})
	bot.Manager().LoadScripts("./scripts")

	reconnect := true
	go func() { <-quit; reconnect = false; bot.Quit() }()

	defer func() {
		bot.Manager().UnloadAll()
		db.Close()
		<-time.After(3 * time.Second)
	}()

	for reconnect {
		if err := bot.Dial(conf.Server, conf.Nick, conf.User); err != nil {
			return
		}

		<-bot.Wait()
		<-time.After(5 * time.Second)
	}
}
