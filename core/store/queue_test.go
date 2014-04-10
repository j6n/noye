package store

import (
	"sync"
	"testing"
	"time"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMessageQueue(t *testing.T) {
	Convey("Given a message queue", t, func() {
		mq := NewQueue()
		Convey("Given a blacklist", func() {
			mq.Blacklist("foo", "bar")
			var res []string

			Convey("It should send to public", func() {
				_, pub := mq.Subscribe("baz", false)
				So(pub, ShouldNotBeNil)

				_, priv := mq.Subscribe("baz", true)
				So(priv, ShouldNotBeNil)

				// send non-blocked as public
				// all listeners should hear it
				res = waitFor(func() { mq.Update("baz", "bar", true) }, pub, priv)
				So(res, ShouldResemble, []string{"bar", "bar"})

				// send non-blocked as private
				// all listeners should hear it
				res = waitFor(func() { mq.Update("baz", "bar", false) }, pub, priv)
				So(res, ShouldResemble, []string{"bar", "bar"})
			})

			Convey("It should send to private", func() {
				_, pub := mq.Subscribe("foo", false)
				So(pub, ShouldBeNil)

				_, priv := mq.Subscribe("foo", true)
				So(priv, ShouldNotBeNil)

				// send blocked as public
				// public shouldn't hear it
				// private should hear it
				res = waitFor(func() { mq.Update("foo", "bar", true) }, pub, priv)
				So(res, ShouldResemble, []string{"?", "bar"})

				// send blocked as private
				// no listeners should hear it
				res = waitFor(func() { mq.Update("foo", "bar", false) }, pub, priv)
				So(res, ShouldResemble, []string{"?", "!"})
			})
		})
	})
}

// this waits on an a list of channels and calls f
// key: ? nil channel, ! timeout, other result
func waitFor(f func(), chs ...chan string) []string {
	var wg sync.WaitGroup
	var out []string

	temp := make(chan string)
	defer close(temp)

	go func() {
		for s := range temp {
			out = append(out, s)
			wg.Done()
		}
	}()

	for _, ch := range chs {
		wg.Add(1)
		if ch == nil {
			temp <- "?"
			continue
		}

		go func(ch chan string) {
			select {
			case s := <-ch:
				temp <- s
			case <-time.After(100 * time.Millisecond):
				temp <- "!"
			}
		}(ch)
	}

	go f()
	wg.Wait()
	return out
}
