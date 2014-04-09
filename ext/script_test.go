package ext

import (
	"fmt"
	"testing"

	"github.com/j6n/noye/mock"
	. "github.com/smartystreets/goconvey/convey"
)

// this should test the script injected methods
// also will serve as script documentation

func TestScriptHTML(t *testing.T) {
	Convey("Given a script", t, func() {
		var done chan struct{}
		bot := mock.NewMockBot()
		bot.SendFn = func(f string, a ...interface{}) {
			t.Logf(f, a...)
			So(fmt.Sprintf(f, a...), ShouldNotBeNil)
			close(done)
		}
		Convey("it should make a parser", func() {
			m := New(bot)
			done = make(chan struct{})
			source := `noye.bot.Send("%+v", html.new("http://example.net/"));`

			m.load(source, "test.js", "/path/to/test.js")
			<-done
		})
		Convey("it should find some stuff", func() {
			m := New(bot)
			done = make(chan struct{})
			source := `
	var doc = html.new("http://music.naver.com/promotion/specialContent.nhn?articleId=4569");
	doc.Find("li[id^='videoContent']", "id")
	noye.bot.Send(">%+v", doc.Results());`

			m.load(source, "test.js", "/path/to/test.js")
			<-done
		})
	})
}
