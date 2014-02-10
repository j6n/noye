package naver

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/j6n/noye/mock"
	"github.com/j6n/noye/noye"
)

func TestNaver(t *testing.T) {
	t.Parallel()

	Convey("Naver plugin should", t, func() {
		Convey("handle a music command", func() {
			SkipConvey("with one url", func() {
				bot, naver, res := mock.NewMockBot(), New(), make(chan string)
				bot.PrivmsgFn = func(target, msg string) { res <- msg }

				naver.Hook(bot)
				naver.Listen() <- noye.Message{
					"museun",
					"#museun",
					"noye: music http://music.naver.com/promotion/specialContent.nhn?articleId=4569&page=2",
				}

				So(<-res, ShouldEndWith, "[720P] 안무 영상")
				So(<-res, ShouldEndWith, "[720P] 포인트 안무(카톡 댄스)")
				So(<-res, ShouldEndWith, "[720P] 인사말 영상")

				close(res)
			})

			SkipConvey("with multiple urls", func() {
				bot, naver, res := mock.NewMockBot(), New(), make(chan string)
				bot.PrivmsgFn = func(target, msg string) { res <- msg }

				naver.Hook(bot)
				naver.Listen() <- noye.Message{
					"museun",
					"#museun",
					"noye: music http://music.naver.com/promotion/specialContent.nhn?articleId=4569&page=2 " +
						"http://music.naver.com/promotion/specialContent.nhn?articleId=4557&page=2",
				}

				So(<-res, ShouldEndWith, "[720P] 안무 영상")
				So(<-res, ShouldEndWith, "[720P] 포인트 안무(카톡 댄스)")
				So(<-res, ShouldEndWith, "[720P] 인사말 영상")
				// --
				So(<-res, ShouldEndWith, "[720P] 메이킹 영상")
				So(<-res, ShouldEndWith, "[720P] 인사말 영상")

				close(res)
			})

			Convey("with mixed tvcast/music", func() {
				bot, naver, res := mock.NewMockBot(), New(), make(chan string)
				bot.PrivmsgFn = func(target, msg string) { res <- msg }

				naver.Hook(bot)
				naver.Listen() <- noye.Message{
					"museun",
					"#museun",
					"noye: music http://music.naver.com/promotion/specialContent.nhn?articleId=4592&page=1",
				}

				t.Log(<-res)
				t.Log(<-res)

				close(res)
			})
		})

		SkipConvey("handle a tvcast command", func() {
			Convey("with one url", func() {
				bot, naver, res := mock.NewMockBot(), New(), make(chan string)
				bot.PrivmsgFn = func(target, msg string) { res <- msg }

				naver.Hook(bot)
				naver.Listen() <- noye.Message{
					"museun",
					"#museun",
					"noye: tvcast http://tvcast.naver.com/v/42788",
				}

				So(<-res, ShouldEndWith, "[720P] [더스타] 걸스데이 소진, 대걸레 잡고 샤이니 '드림 걸' 안무시범 '폭소'")
				close(res)
			})

			Convey("with one id", func() {
				bot, naver, res := mock.NewMockBot(), New(), make(chan string)
				bot.PrivmsgFn = func(target, msg string) { res <- msg }

				naver.Hook(bot)
				naver.Listen() <- noye.Message{
					"museun",
					"#museun",
					"noye: tvcast 42788",
				}

				So(<-res, ShouldEndWith, "[720P] [더스타] 걸스데이 소진, 대걸레 잡고 샤이니 '드림 걸' 안무시범 '폭소'")
				close(res)
			})

			Convey("with multiple urls", func() {
				bot, naver, res := mock.NewMockBot(), New(), make(chan string)
				bot.PrivmsgFn = func(target, msg string) { res <- msg }

				naver.Hook(bot)
				naver.Listen() <- noye.Message{
					"museun",
					"#museun",
					"noye: tvcast http://tvcast.naver.com/v/42788 http://tvcast.naver.com/v/42782",
				}

				So(<-res, ShouldEndWith, "[720P] [더스타] 걸스데이 소진, 대걸레 잡고 샤이니 '드림 걸' 안무시범 '폭소'")
				So(<-res, ShouldEndWith, "[720P] [더스타] 걸스데이 민아, 3콤보 라이브+웨이브+애교…'방만능'된 사연")
				close(res)
			})

			Convey("with urls and ids", func() {
				bot, naver, res := mock.NewMockBot(), New(), make(chan string)
				bot.PrivmsgFn = func(target, msg string) { res <- msg }

				naver.Hook(bot)
				naver.Listen() <- noye.Message{
					"museun",
					"#museun",
					"noye: tvcast http://tvcast.naver.com/v/42788 42782",
				}

				So(<-res, ShouldEndWith, "[720P] [더스타] 걸스데이 소진, 대걸레 잡고 샤이니 '드림 걸' 안무시범 '폭소'")
				So(<-res, ShouldEndWith, "[720P] [더스타] 걸스데이 민아, 3콤보 라이브+웨이브+애교…'방만능'된 사연")
				close(res)
			})
		})
	})
}
