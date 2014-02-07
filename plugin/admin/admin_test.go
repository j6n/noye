package admin

import (
	"testing"

	"github.com/j6n/noye/mock"
	"github.com/j6n/noye/noye"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAdmin(t *testing.T) {
	t.Parallel()

	Convey("Admin plugin should", t, func() {
		Convey("Be able to be created", func() {
			admin := New()
			So(admin, ShouldNotBeNil)
		})

		Convey("Handle a join command", func() {
			bot, admin, res := mock.NewMockBot(), New(), make(chan string)
			bot.JoinFn = func(target string) { res <- target; close(res) }

			admin.base.Hook(bot)
			admin.base.Messages <- noye.Message{"museun", "#museun", "noye: join #foobar"}

			So(<-res, ShouldEqual, "#foobar")
		})

		Convey("Handle a part command", func() {
			bot, admin, res := mock.NewMockBot(), New(), make(chan string)
			bot.PartFn = func(target string) { res <- target; close(res) }

			admin.base.Hook(bot)
			admin.base.Messages <- noye.Message{"museun", "#museun", "noye: part #foobar"}

			So(<-res, ShouldEqual, "#foobar")
		})
	})
}
