package admin

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/j6n/noye/mock"
	"github.com/j6n/noye/noye"
)

func TestAdmin(t *testing.T) {
	t.Parallel()

	Convey("Admin plugin should", t, func() {
		Convey("Handle a join command", func() {
			bot, admin, res := mock.NewMockBot(), New(), make(chan string)
			bot.JoinFn = func(target string) { res <- target; close(res) }

			admin.Hook(bot)
			admin.Listen() <- noye.Message{"museun", "#museun", "noye: join #foobar"}

			So(<-res, ShouldEqual, "#foobar")
		})

		Convey("Handle a join command, multiple channels", func() {
			bot, admin, res := mock.NewMockBot(), New(), make(chan string)
			bot.JoinFn = func(target string) { res <- target }

			admin.Hook(bot)
			admin.Listen() <- noye.Message{"museun", "#museun", "noye: join #foobar #testing"}

			So(<-res, ShouldEqual, "#foobar")
			So(<-res, ShouldEqual, "#testing")
			close(res)
		})

		Convey("Handle a part command", func() {
			bot, admin, res := mock.NewMockBot(), New(), make(chan string)
			bot.PartFn = func(target string) { res <- target; close(res) }

			admin.Hook(bot)
			admin.Listen() <- noye.Message{"museun", "#museun", "noye: part #foobar"}

			So(<-res, ShouldEqual, "#foobar")
		})

		Convey("Handle a part command, multiple channels", func() {
			bot, admin, res := mock.NewMockBot(), New(), make(chan string)
			bot.PartFn = func(target string) { res <- target }

			admin.Hook(bot)
			admin.Listen() <- noye.Message{"museun", "#museun", "noye: part #foobar #testing"}

			So(<-res, ShouldEqual, "#foobar")
			So(<-res, ShouldEqual, "#testing")
			close(res)
		})
	})
}
