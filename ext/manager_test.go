package ext

import (
	"fmt"
	"testing"
	"time"

	"github.com/j6n/noye/mock"
	"github.com/j6n/noye/noye"
	. "github.com/smartystreets/goconvey/convey"
)

func TestManager(t *testing.T) {
	// quick fail
	fail := func(f string, a ...interface{}) {
		t.Errorf(f, a...)
		t.FailNow()
	}

	Convey("Given a manager", t, func() {
		// create a mock bot
		ctx := mock.NewMockBot()

		// crate a manager
		manager := New(ctx)

		// this returns and logs a privmsg received
		privmsg := func() <-chan string {
			out := make(chan string)
			ctx.PrivmsgFn = func(target, msg string) {
				res := fmt.Sprintf("%s: %s", target, msg)
				t.Log("PRIVMSG >>", res)
				out <- res
			}
			return out
		}

		// easily make a new script from source
		script := func(source string) error {
			return manager.load(source, "/this/test/script.js")
		}

		// this creates a new noye.Message and sends it to the manager
		respond := func(text string, other ...string) {
			target, from := "#noye", "#museun"
			if len(other) == 2 {
				target, from = other[0], other[1]
			}
			manager.Respond(noye.Message{Text: text, Target: target, From: from})
		}

		Convey("it should respond", func() {
			err := script(`
			respond("!(hello|bye)", function(msg) {
				msg.Reply("hello!");
			});`)
			So(err, ShouldBeNil)

			out := privmsg()
			respond("!hello test")

			select {
			case msg := <-out:
				t.Log("out:", msg)
			case <-time.After(3 * time.Second):
				fail("timed out")
			}
		})
	})
}
