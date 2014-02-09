package plugin

import (
	"regexp"
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
			cmd := &Command{Command: "hello"}
			So(match(cmd, "hello"), ShouldBeTrue)
			So(match(cmd, "something"), ShouldBeFalse)
		})

		Convey("should match a simple respond", func() {
			cmd := &Command{Respond: true, Command: "something"}
			So(match(cmd, "noye: something"), ShouldBeTrue)
			So(match(cmd, "noye: do something"), ShouldBeFalse)
		})

		Convey("should match multiple parts", func() {
			cmd := &Command{Command: "foo", Each: true,
				Matcher: func(s string) (bool, string) { return len(s) == 3, "" },
			}
			So(match(cmd, "foo bar baz"), ShouldBeTrue)
			So(match(cmd, "foo foo foobar"), ShouldBeFalse)
			So(match(cmd, "noye: test this out"), ShouldBeFalse)
		})

		Convey("should match respond with mulitple parts", func() {
			cmd := &Command{Command: "foo", Each: true, Respond: true,
				Matcher: func(s string) (bool, string) { return len(s) == 3, "" },
			}
			So(match(cmd, "noye: foo bar baz"), ShouldBeTrue)
			So(match(cmd, "noye: foo bar asdf"), ShouldBeFalse)
			So(match(cmd, "foo bar asdf"), ShouldBeFalse)
			So(match(cmd, "foo bar baz"), ShouldBeFalse)
			So(match(cmd, "noye: bar foo"), ShouldBeFalse)
		})

		Convey("should match simple with a result", func() {
			cmd := &Command{Command: "foo", Matcher: func(s string) (bool, string) {
				t.Log(">", s)
				if s == "test" {
					return true, "bar"
				}
				return false, ""
			}}

			So(match(cmd, "foo test"), ShouldBeTrue)

			res := cmd.Results()
			So(res, ShouldNotBeNil)
			So(len(res), ShouldEqual, 1)
			So(res[0], ShouldEqual, "bar")
		})

		Convey("should match multiple with results", func() {
			cmd := &Command{Command: "foo", Each: true, Matcher: func(s string) (bool, string) {
				ok, _ := regexp.MatchString("[0-9]", s)
				if ok {
					return ok, s
				}

				return false, ""
			}}

			So(match(cmd, "foo 1 0 0 4"), ShouldBeTrue)

			res := cmd.Results()
			So(res, ShouldNotBeNil)
			So(len(res), ShouldEqual, 4)
			So(res, ShouldResemble, []string{"1", "0", "0", "4"})
		})
	})
}
