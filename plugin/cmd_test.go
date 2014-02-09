package plugin

import (
	"strings"
	"testing"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/j6n/noye/noye"
)

func TestCommand(t *testing.T) {
	// helper for less typing
	// it builds a simple message and matches it to the command
	match := func(cmd *Command, s ...string) bool {
		return cmd.Match(noye.Message{"museun", "#museun", strings.Join(s, " ")})
	}

	Convey("Command", t, func() {
		Convey("should match a simple command", func() {
			cmd := Command{Command: "hello"}
			So(match(&cmd, "hello"), ShouldBeTrue)
			So(match(&cmd, "something"), ShouldBeFalse)
		})

		Convey("should match a simple respond", func() {
			cmd := Command{Respond: true, Command: "something"}
			So(match(&cmd, "noye: something"), ShouldBeTrue)
			So(match(&cmd, "noye: do something"), ShouldBeFalse)
		})

		Convey("should match multiple parts", func() {
			cmd := Command{Command: "foo", Each: true,
				Matcher: func(s string) bool { return len(s) == 3 },
			}
			So(match(&cmd, "foo bar baz"), ShouldBeTrue)
			So(match(&cmd, "foo foo foobar"), ShouldBeFalse)
			So(match(&cmd, "noye: test this out"), ShouldBeFalse)
		})

		Convey("should match respond with mulitple parts", func() {
			cmd := Command{Command: "foo", Each: true, Respond: true,
				Matcher: func(s string) bool { return len(s) == 3 },
			}
			So(match(&cmd, "noye: foo bar baz"), ShouldBeTrue)
			So(match(&cmd, "noye: foo bar asdf"), ShouldBeFalse)
			So(match(&cmd, "foo bar asdf"), ShouldBeFalse)
			So(match(&cmd, "foo bar baz"), ShouldBeFalse)
			So(match(&cmd, "noye: bar foo"), ShouldBeFalse)
		})
	})
}
