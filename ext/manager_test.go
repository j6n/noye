package ext

import (
	"testing"

	"github.com/j6n/logger"
	"github.com/j6n/noye/mock"
	"github.com/j6n/noye/noye"
	. "github.com/smartystreets/goconvey/convey"
)

func TestManager(t *testing.T) {
	log.Formatter = &logger.JsonFormatter{true}

	ctx := mock.NewMockBot()
	ctx.PrivmsgFn = func(target, msg string) {
		t.Logf("PRIVMSG >> %s: %s\n", target, msg)
	}
	manager := New(ctx)

	Convey("Given a manager", t, func() {
		source := `
		respond("!(hello|bye)", function(msg) {
			log("hello or bye " + msg.Text);
			log("from: " + msg.From + " on " + msg.Target);
			noye.reply(msg, "hello!");
		});
		`
		path := "/this/test/script.js"
		err := manager.load(source, path)
		So(err, ShouldBeNil)

		manager.Respond(noye.Message{
			From:   "museun",
			Target: "#noye",
			Text:   "!hello test",
		})
	})
}
