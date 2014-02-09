package naver

import (
	"regexp"

	"github.com/j6n/noye-naver"
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
	//music := dsl.Nick("noye").Command("naver").Param("music").List(`(http://music.naver.com/.*?\S*)+`)

	music := plugin.Command{
		Respond: true,
		Command: "naver",
		Each:    true,
		Matcher: plugin.RegexMatcher(
			regexp.MustCompile(`(http://music.naver.com/.*?\S*)+`),
			true,
		),
	}

	for msg := range n.Listen() {
		switch {
		case music.Match(msg):
			n.handleMusic(msg, music.Results())
		}
	}
}

func (n *Naver) handleMusic(msg noye.Message, match []string) {
	for _, url := range match {
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
