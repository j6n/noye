package ext

import (
	"testing"

	"github.com/j6n/noye/mock"
	"github.com/j6n/noye/noye"
	. "github.com/smartystreets/goconvey/convey"
)

func TestManager(t *testing.T) {
	ctx := mock.NewMockBot()
	manager := New(ctx)
	ctx.PrivmsgFn = func(target, msg string) {
		t.Logf("PRIVMSG >> %s: %s\n", target, msg)
	}

	Convey("Given a manager", t, func() {
		Convey("should respond", func() {
			source := `
			respond("!(hello|bye)", function(msg) {
				log("hello or bye " + msg.Text);
				log("from: " + msg.From + " on " + msg.Target);
				noye.reply(msg, "hello!");
			});`

			path := "/this/test/script.js"
			err := manager.load(source, path)
			So(err, ShouldBeNil)

			manager.Respond(noye.Message{
				From:   "museun",
				Target: "#noye",
				Text:   "!hello test",
			})
		})

		Convey("should respond and have results", func() {
			source := `
			respond("!foo (bar|baz)$", function(msg, res) {
				log("foo " + res[1]);
			});`

			path := "/this/test/script.js"
			err := manager.load(source, path)
			So(err, ShouldBeNil)

			manager.Respond(noye.Message{
				From:   "museun",
				Target: "#noye",
				Text:   "!foo bar",
			})
		})
	})

	source := `
	respond("!ip", function(msg) {
		var data = core.http("http://ifconfig.me/ip").get();
		noye.reply(msg, data);
	});`

	path := "/this/test/script.js"
	err := manager.load(source, path)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	manager.Respond(noye.Message{
		From:   "museun",
		Target: "#noye",
		Text:   "!ip",
	})
}
