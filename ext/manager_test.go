package ext

//import . "github.com/smartystreets/goconvey/convey"

// func TestManager(t *testing.T) {
// 	fail := func(f string, a ...interface{}) {
// 		t.Errorf(f, a...)
// 		//t.FailNow()
// 	}

// 	Convey("Given a manager", t, func() {
// 		ctx := mock.NewMockBot()
// 		manager := New(ctx)

// 		script := func(source string) error { return manager.load(source, "/this/test/script.js") }
// 		listen := func(ev string) { manager.Listen(noye.IrcMessage{Command: ev, Source: "museun!much@localhost"}) }

// 		respond := func(text string, other ...string) {
// 			from, target := "museun", "#noye"
// 			if len(other) == 2 {
// 				from, target = other[0], other[1]
// 			}

// 			manager.Respond(noye.Message{from, target, text})
// 		}

// 		Convey("it should respond", func() {
// 			err := script(`
// 			respond("!hello(?:$|\s*(?P<out>.+$))", function(msg, res) {
// 				if (res.out) {
// 					return msg.Send("hello %s", res.out);
// 				}

// 				msg.Reply("hello!");
// 			});`)
// 			So(err, ShouldBeNil)

// 			res := make(chan string)
// 			ctx.PrivmsgFn = func(target, msg string) {
// 				out := fmt.Sprintf("%s: %s", target, msg)
// 				res <- out
// 			}

// 			respond("!hello")
// 			select {
// 			case input := <-res:
// 				So(input, ShouldEqual, "#noye: museun: hello!")
// 			case <-time.After(3 * time.Second):
// 				fail("timed out")
// 			}

// 			respond("!hello test")
// 			select {
// 			case input := <-res:
// 				So(input, ShouldEqual, "#noye: hello test")
// 			case <-time.After(3 * time.Second):
// 				fail("timed out")
// 			}

// 		})

// 		Convey("it should listen", func() {
// 			err := script(`
// 			listen("001", function(msg) {
// 				noye.Join("#test");
// 			});`)
// 			So(err, ShouldBeNil)

// 			res := make(chan string)
// 			ctx.JoinFn = func(target string) {
// 				res <- target
// 			}

// 			listen("001")
// 			select {
// 			case input := <-res:
// 				So(input, ShouldEqual, "#test")
// 				return
// 			case <-time.After(3 * time.Second):
// 				fail("timed out")
// 			}
// 		})
// 	})
// }
