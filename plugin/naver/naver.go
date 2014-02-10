package naver

import (
	"regexp"

	"github.com/j6n/naver/music"
	"github.com/j6n/naver/tvcast"
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

var tvcastRe = regexp.MustCompile(`(?:tvcast.naver.com/v/(\d+)$)|(\d+)$`)
var musicRe = regexp.MustCompile(`(http://music.naver.com/.*?\S*)+`)

func (n *Naver) process() {
	music := plugin.Respond("music", plugin.RegexMatcher(musicRe, true))
	music.Each = true

	tvcast := plugin.Respond("tvcast", plugin.RegexMatcher(tvcastRe, true))
	tvcast.Each = true

	for msg := range n.Listen() {
		switch {
		case music.Match(msg):
			n.handleMusic(msg, music.Results())
		case tvcast.Match(msg):
			n.handleTvcast(msg, tvcast.Results())
		}
	}
}

func (n *Naver) handleMusic(msg noye.Message, match []string) {
	defer recover() // don't crash

	for _, url := range match {
		ids, err := music.FindIDs(url)
		if err != nil {
			n.Error(msg, "music/findIDs", err)
			continue
		}

		for _, id := range ids {
			vid, err := music.GetVideo(id)
			if err != nil {
				n.Error(msg, "music: "+id, err)
				continue
			}

			n.Reply(msg, "[%s] %s", vid.Encoding, vid.Title)
			n.Reply(msg, "%s", vid.PlayUrl)
		}
	}
}

func (n *Naver) handleTvcast(msg noye.Message, matches []string) {
	defer recover() // don't crash

	for _, match := range matches {
		var id string
		for _, p := range tvcastRe.FindStringSubmatch(match)[1:] {
			if p != "" {
				id = p
				break
			}
		}

		vid, err := tvcast.GetVideo(id)
		if err != nil {
			n.Error(msg, "tvcast: "+id, err)
			continue
		}

		n.Reply(msg, "[%s] %s", vid.Encoding, vid.Title)
		n.Reply(msg, "%s", vid.PlayUrl)
	}
}
