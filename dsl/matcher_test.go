package dsl

import (
	"regexp"
	"testing"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/j6n/noye/noye"
)

func TestMatcher(t *testing.T) {
	t.Parallel()

	Convey("Matcher should", t, func() {
		Convey("Match command with a list", func() {
			cmd := New().Command("foo|bar").List("([0-9])")
			ok, err := cmd.Valid()
			So(ok, ShouldBeTrue)
			So(err, ShouldBeNil)

			ok = cmd.Match(noye.Message{"museun", "#museun", "foo 1 3 3 7"})
			So(ok, ShouldBeTrue)
			So(cmd.Results.Cmds(), ShouldBeZeroValue)
			So(cmd.Results.Params(), ShouldBeZeroValue)
			So(cmd.Results.Lists(), ShouldResemble, []string{"1", "3", "3", "7"})

			ok = cmd.Match(noye.Message{"museun", "#museun", "bar 1 0 0 4"})
			So(ok, ShouldBeTrue)
			So(cmd.Results.Cmds(), ShouldBeZeroValue)
			So(cmd.Results.Params(), ShouldBeZeroValue)
			So(cmd.Results.Lists(), ShouldResemble, []string{"1", "0", "0", "4"})
		})

		Convey("Match command with some params", func() {
			cmd := New().Command("foo").Param("([0-9])").Param("([a-z])")
			ok, err := cmd.Valid()
			So(ok, ShouldBeTrue)
			So(err, ShouldBeNil)

			ok = cmd.Match(noye.Message{"museun", "#museun", "foo 7 z"})
			So(ok, ShouldBeTrue)
			So(cmd.Results.Cmds(), ShouldBeZeroValue)
			So(cmd.Results.Params(), ShouldResemble, []string{"7", "z"})
			So(cmd.Results.Lists(), ShouldBeZeroValue)

			ok = cmd.Match(noye.Message{"museun", "#museun", "foo 9 9"})
			So(ok, ShouldBeFalse)

			ok = cmd.Match(noye.Message{"museun", "#museun", "bar 1 0 0 4"})
			So(ok, ShouldBeFalse)
		})

		Convey("Match nick prefix and command", func() {
			cmd := Nick("noye").Command("foo")
			ok, err := cmd.Valid()
			So(ok, ShouldBeTrue)
			So(err, ShouldBeNil)

			ok = cmd.Match(noye.Message{"museun", "#museun", "noye foo bar baz"})
			So(ok, ShouldBeTrue)

			ok = cmd.Match(noye.Message{"museun", "#museun", "foo 9 9"})
			So(ok, ShouldBeFalse)

			ok = cmd.Match(noye.Message{"museun", "#museun", "noye baz foo"})
			So(ok, ShouldBeFalse)

			So(cmd.Results.Cmds(), ShouldBeZeroValue)
			So(cmd.Results.Params(), ShouldBeZeroValue)
			So(cmd.Results.Lists(), ShouldBeZeroValue)
		})

		Convey("Match nick prefix, command and param", func() {
			cmd := Nick("noye").Command("(foo)").Param("([0-9])")
			ok, err := cmd.Valid()
			So(ok, ShouldBeTrue)
			So(err, ShouldBeNil)

			ok = cmd.Match(noye.Message{"museun", "#museun", "noye foo 7 bar"})
			So(ok, ShouldBeTrue)
			So(cmd.Results.Cmds(), ShouldResemble, []string{"foo"})
			So(cmd.Results.Params(), ShouldResemble, []string{"7"})

			ok = cmd.Match(noye.Message{"museun", "#museun", "foo 9 9"})
			So(ok, ShouldBeFalse)

			ok = cmd.Match(noye.Message{"museun", "#museun", "noye foo 9 baz"})
			So(ok, ShouldBeTrue)

			So(cmd.Results.Cmds(), ShouldResemble, []string{"foo"})
			So(cmd.Results.Params(), ShouldResemble, []string{"9"})
			So(cmd.Results.Lists(), ShouldBeZeroValue)
		})
	})
}

func TestNickRegex(t *testing.T) {
	Convey("Regex should match", t, func() {
		nick := "noye"
		re, err := regexp.Compile(`(?:` + nick + `[:,]?\s*)`)
		So(err, ShouldBeNil)

		So(re.MatchString("noye"), ShouldBeTrue)
		So(re.MatchString("noye:"), ShouldBeTrue)
		So(re.MatchString("noye "), ShouldBeTrue)
		So(re.MatchString("noye: "), ShouldBeTrue)
	})
}
