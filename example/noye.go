package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/j6n/noye-naver"
	"github.com/j6n/noye/irc"
)

func main() {
	log.Println("Starting connection")
	conn := &irc.Connection{}

	bot := irc.New(conn)
	bot.Autojoin = []string{"#museun", "#nanashin"}
	bot.Handle = func(msg irc.Message) {
		if matches := re.FindAllStringSubmatch(msg.Text, -1); len(matches) > 0 && len(matches[0]) > 0 {
			for _, match := range matches[0][1:] {
				ids, err := naver.FindIDs(match)
				if err != nil {
					log.Println(err)
					continue
				}
				for _, id := range ids {
					video, err := naver.GetVideo(id)
					if err != nil {
						log.Println(err)
						continue
					}
					bot.Send("PRIVMSG %s :[%s] %s", msg.Args[0], video.Encoding, video.Title)
					bot.Send("PRIVMSG %s :%s", msg.Args[0], video.PlayUrl)
				}
			}
		}
	}

	if err := bot.Dial("irc.quakenet.org:6667", "noye", "museun"); err != nil {
		log.Fatalln(err)
	}

	<-bot.Wait()
	log.Println("done")
}

var re = regexp.MustCompile(strings.Join([]string{
	"(?i)",
	"(?:\\A",
	"(?:" + "noye" + "[:,]?\\s*|/)",
	"(?:" + "naver" + "\\s*)|\\S*)",
	"(?:" + `(http://music.naver.com/.*?\S*)+` + ")",
	"(?:\\S*|\\z)",
}, ""))
