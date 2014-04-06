package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"runtime"

	"github.com/j6n/noye/store"

	"github.com/j6n/noye/irc"
	"github.com/j6n/noye/logger"
)

var (
	log  = logger.Get()
	conf = NewConfig()
	db   *store.DB
)

func init() {
	runtime.GOMAXPROCS(4)

	db, _ = store.NewDB()
	if err := db.CheckTable("config", store.KvSchema); err != nil {
		log.Fatalf("can't create table %s:%s\n", "config", err)
	}

	m := conf.toMap()
	for k, v := range m {
		db.Set("config", k, v)
	}
}

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	bot := irc.New(&irc.Connection{})
	ext := bot.Manager()

	scripts := getFiles("./scripts")
	for script := range scripts {
		log.Infof("found script: '%s'\n", script)
		if err := ext.Load(script); err != nil {
			log.Errorf("loading script '%s': '%s'\n", script, err)
		}
	}

	reconnect := true
	go func() { <-quit; reconnect = false; bot.Quit() }()

	for reconnect {
		if err := bot.Dial(conf.Server, conf.Nick, conf.Nick); err != nil {
			return
		}

		<-bot.Wait()
	}
}

func getFiles(base string) <-chan string {
	scripts := make(chan string)
	go func() {
		walker := func(fp string, fi os.FileInfo, err error) error {
			if err != nil || !!fi.IsDir() {
				return nil
			}
			matched, err := filepath.Match("*.js", fi.Name())
			if err != nil {
				return err
			}
			if matched {
				scripts <- fp
			}
			return nil
		}

		if err := filepath.Walk(base, walker); err != nil {
			log.Errorf("Walking '%s': %s\n", base, err)
		}
		close(scripts)
	}()

	return scripts
}
