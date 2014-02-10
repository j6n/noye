package urlinfo

import (
	"net/http"
	"testing"
	"time"

	"github.com/j6n/noye/mock"
	"github.com/j6n/noye/noye"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUrlinfo(t *testing.T) {
	Convey("Urlinfo should", t, func() {
		bot, urlinfo := mock.NewMockBot(), New()
		//bot.PrivmsgFn = func(target, msg string) { res <- msg }

		urlinfo.Hook(bot)
		urlinfo.Listen() <- noye.Message{
			"museun",
			"#museun",
			"http://google.com http://www.youtube.com/watch?v=FkwrdyWzyDc http://i.imgur.com/vV4mtFs.gif " +
				"https://pbs.twimg.com/media/BT4soiDIYAA6KlT.png:large " +
				"http://cfile8.uf.tistory.com/original/271B9F3B5228BEC234223D " +
				"http://download.jetbrains.com/resharper/ReSharperSetup.8.0.1000.2286.msi",
		}

		<-time.After(10 * time.Second)
	})
}

func _testHead(t *testing.T) {
	resp, err := http.Head("http://google.com")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(resp.Status)
	t.Logf("%+v\n", resp.Header)
	t.Log(resp.ContentLength)
}

func _testRegex(t *testing.T) {
	input := []string{
		"http://google.com",
		"http://www.youtube.com/watch?v=FkwrdyWzyDc",
		"http://i.imgur.com/vV4mtFs.gif",
		"https://pbs.twimg.com/media/BT4soiDIYAA6KlT.png:large",
		"http://cfile8.uf.tistory.com/original/271B9F3B5228BEC234223D",
		"http://download.jetbrains.com/resharper/ReSharperSetup.8.0.1000.2286.msi",
	}

	for _, url := range input {
		t.Log(urlinfoRe.MatchString(url), url)
		t.Log(urlinfoRe.FindString(url))
		t.Log("------")
	}
}
