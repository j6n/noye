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
			mq.Blacklist("foo", "bar")

			Convey("It should send to public", func() {
				_, ch := mq.Subscribe("baz")
				So(ch, ShouldNotBeNil)

				var res string
				go func() {
					select {
					case res = <-ch:
					}
				}()

				// send non-blocked as private
				// all listeners should hear it
				mq.Update("baz", "bar", true)
				So(res, ShouldEqual, "bar")

				res = ""
				go func() {
					select {
					case res = <-ch:
					case <-time.After(100 * time.Millisecond):
					}
				}()

				// send non-blocked as public
				// all listeners should hear it
				mq.Update("baz", "bar", false)
				So(res, ShouldEqual, "bar")
			})

			Convey("It should send to private", func() {
				_, ch := mq.Subscribe("foo")
				So(ch, ShouldNotBeNil)

				var res string
				go func() {
					select {
					case res = <-ch:
					}
				}()

				// send blocked as public
				// all listeners should hear it
				mq.Update("foo", "bar", false)
				So(res, ShouldEqual, "bar")

				res = ""
				go func() {
					select {
					case res = <-ch:
					case <-time.After(100 * time.Millisecond):
					}
				}()

				// send blocked as private
				// no listeners should hear it
				mq.Update("foo", "bar", true)
				So(res, ShouldBeEmpty)
			})
		})
	})
}
