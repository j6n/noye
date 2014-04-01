package main

import (
	"os"
	"os/signal"
	"path/filepath"
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

	bot := irc.New(&irc.Connection{})
	ext := bot.Manager()

	scripts := getFiles("./scripts")
	for script := range scripts {
		if err := ext.Load(script); err != nil {
			// log error
			err = nil
		}
	}

	if err := bot.Dial("irc.quakenet.org:6667", "noye", "museun"); err != nil {
		return
	}

	go func() {
		// wait for the ctrl-c
		<-quit
		// send the quit
		bot.Quit()
	}()

	// wait for the close signal
	<-bot.Wait()
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

		filepath.Walk(base, walker)
		close(scripts)
	}()

	return scripts
}
