package store

import (
	"testing"
	"time"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMessageQueue(t *testing.T) {
	Convey("Given a message queue", t, func() {
		mq := NewQueue()
		Convey("Given a blacklist", func() {
			mq.Blacklist("foo", "bar", "sojin")

			Convey("It should ignore matching public requests", func() {
				id, ch := mq.Subscribe("foo", false)
				So(id, ShouldEqual, 0)
				So(ch, ShouldBeNil)
				mq.Update("foo", "bar", true)
				mq.Update("foo", "bar", false)
			})

			Convey("It should honor matching private requests", func() {
				_, ch := mq.Subscribe("foo", true)
				So(ch, ShouldNotBeNil)

				var res string
				go func() {
					select {
					case res = <-ch:
					}
				}()

				mq.Update("foo", "bar", true)
				So(res, ShouldEqual, "bar")

				res = ""
				go func() {
					select {
					case res = <-ch:
					case <-time.After(100 * time.Millisecond):
					}
				}()

				mq.Update("foo", "bar", false)
				So(res, ShouldBeEmpty)
			})
		})
	})
}
