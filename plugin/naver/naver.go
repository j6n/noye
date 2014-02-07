package naver

import (
	"log"

	"github.com/j6n/noye-naver"

	"github.com/j6n/noye/dsl"
	"github.com/j6n/noye/noye"
	"github.com/j6n/noye/plugin"
)

type Naver struct {
	*plugin.BasePlugin
}

func New() *Naver {
	naver := &Naver{plugin.New()}
	go naver.process()
	return naver
}

func (n *Naver) process() {
	music := dsl.Nick("noye").Command("naver").Param("music").List(`(http://music.naver.com/.*?\S*)+`)

	if ok, err := music.Valid(); !ok {
		log.Println("err starting naver:", err)
		return
	}

	for msg := range n.Messages {
		switch {
		case music.Match(msg):
			n.handleMusic(msg, &music.Results)
		}
	}
}

func (n *Naver) handleMusic(msg noye.Message, match *dsl.Results) {
	for _, url := range match.Lists() {
		ids, err := naver.FindIDs(url)
		if err != nil {
			n.Error(msg, "music/findIDs", err)
			continue
		}

		for _, id := range ids {
			vid, err := naver.GetVideo(id)
			if err != nil {
				n.Error(msg, "music/findIDs", err)
				continue
			}

			n.Reply(msg, "[%s] %s", vid.Encoding, vid.Title)
			n.Reply(msg, "%s", vid.PlayUrl)
		}
	}
}
