package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"runtime"

	"github.com/j6n/noye/irc"
	"github.com/j6n/noye/logger"
)

var log = logger.Get()
var server string

func init() {
	runtime.GOMAXPROCS(4)

	if server = os.Getenv("NOYE_SERVER"); server == "" {
		server = "localhost:6667"
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
	go func() { <-quit; bot.Quit(); reconnect = false }()

	for reconnect {
		if err := bot.Dial(server, "noye", "museun"); err != nil {
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
