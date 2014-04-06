package irc

import (
	"fmt"
	"testing"

	"github.com/j6n/noye/mock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBot(t *testing.T) {
	Convey("Given a bot", t, func() {
		var done int

		conn := mock.NewMockConn()
		conn.CloseFn = func() { done++ }
		conn.ReadLineFn = func() (string, error) {
			return "NOTICE AUTH :*** Checking Ident", fmt.Errorf("closing")
		}

		test := New(conn)
		Convey("it should reconnect", func() {
			for i := 0; i < 3; i++ {
				err := test.Dial("localhost:6667", "noye", "test")
				So(err, ShouldBeNil)
				<-test.Wait()
			}

			So(done, ShouldEqual, 3)
		})
	})
}
