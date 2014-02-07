package naver

import (
	"testing"

	"github.com/j6n/noye/mock"
	"github.com/j6n/noye/noye"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNaver(t *testing.T) {
	t.Parallel()

	Convey("Naver plugin should", t, func() {
		Convey("Be able to be created", func() {
			naver := New()
			So(naver, ShouldNotBeNil)
		})

		Convey("Handle a music command", func() {
			Convey("With one url", func() {
				bot, naver, res := mock.NewMockBot(), New(), make(chan string)
				bot.PrivmsgFn = func(target, msg string) { res <- msg }

				naver.Hook(bot)
				naver.Messages <- noye.Message{
					"museun",
					"#museun",
					"noye: naver music http://music.naver.com/promotion/specialContent.nhn?articleId=4569&page=2",
				}

				So(<-res, ShouldEqual, "[720P] 안무 영상")
				So(<-res, ShouldNotBeNil) // the url
				So(<-res, ShouldEqual, "[720P] 포인트 안무(카톡 댄스)")
				So(<-res, ShouldNotBeNil) // the url
				So(<-res, ShouldEqual, "[720P] 인사말 영상")
				So(<-res, ShouldNotBeNil) // the url

				close(res)
			})

			FocusConvey("With multiple urls", func() {
				bot, naver, res := mock.NewMockBot(), New(), make(chan string)
				bot.PrivmsgFn = func(target, msg string) { res <- msg }

				naver.Hook(bot)
				naver.Messages <- noye.Message{
					"museun",
					"#museun",
					"noye: naver music http://music.naver.com/promotion/specialContent.nhn?articleId=4569&page=2 " +
						"http://music.naver.com/promotion/specialContent.nhn?articleId=4557&page=2",
				}

				So(<-res, ShouldEqual, "[720P] 안무 영상")
				So(<-res, ShouldNotBeNil) // the url
				So(<-res, ShouldEqual, "[720P] 포인트 안무(카톡 댄스)")
				So(<-res, ShouldNotBeNil) // the url
				So(<-res, ShouldEqual, "[720P] 인사말 영상")
				So(<-res, ShouldNotBeNil) // the url
				// --
				So(<-res, ShouldEqual, "[720P] 메이킹 영상")
				So(<-res, ShouldNotBeNil) // the url
				So(<-res, ShouldEqual, "[720P] 인사말 영상")
				So(<-res, ShouldNotBeNil) // the url

				close(res)
			})
		})
	})
}
