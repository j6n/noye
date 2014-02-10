package translate

import (
	"testing"

	"github.com/j6n/noye/mock"
	"github.com/j6n/noye/noye"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTranslate(t *testing.T) {
	Convey("Translate should", t, func() {
		bot, translate, res := mock.NewMockBot(), New(), make(chan string)
		bot.PrivmsgFn = func(target, msg string) { res <- msg }

		translate.Hook(bot)
		translate.Listen() <- noye.Message{
			"museun",
			"#museun",
			"noye: translate en,ko hello world",
		}

		text := <-res
		So(text, ShouldEqual, "안녕하세요!")

		translate.Listen() <- noye.Message{
			"museun",
			"#museun",
			"noye: translate ko,en " + text,
		}

		So(<-res, ShouldEqual, "Hello there!")
	})
}
