package ext

import (
	"testing"
	"time"

	"github.com/j6n/logger"
	"github.com/j6n/noye/mock"
	"github.com/j6n/noye/noye"
	"github.com/robertkrimen/otto"
	. "github.com/smartystreets/goconvey/convey"
)

func TestManager(t *testing.T) {
	ctx := mock.NewMockBot()
	manager := New(ctx, logger.New())
	log.Formatter = &logger.JsonFormatter{true}
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

	<-time.After(5 * time.Second)
}

func TestFoo(t *testing.T) {
	vm := otto.New()
	if _, err := vm.Run(`console.log(encodeURI("http://ifconfig.me/ip"))`); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
