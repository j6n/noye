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

	newCommand := func(cmd string, matchers ...Matcher) *Command {
		return &Command{Command: cmd, Matchers: matchers}
	}

	Convey("Command should", t, func() {
		Convey("match a simple command", func() {
			cmd := newCommand("hello")

			So(match(cmd, "hello"), ShouldBeTrue)
			So(match(cmd, "something"), ShouldBeFalse)
		})

		Convey("match a simple respond", func() {
			cmd := newCommand("something")
			cmd.Options.Respond = true

			So(match(cmd, "noye: something"), ShouldBeTrue)
			So(match(cmd, "noye: do something"), ShouldBeFalse)
		})

		Convey("match multiple parts", func() {
			cmd := newCommand("foo", func(s string) (bool, string) { return len(s) == 3, "" })
			cmd.Options.Each = true

			So(match(cmd, "foo bar baz"), ShouldBeTrue)
			So(match(cmd, "foo foo true"), ShouldBeTrue) // non strict
			So(match(cmd, "noye: test this out"), ShouldBeFalse)
		})

		Convey("match multiple parts, strict", func() {
			cmd := newCommand("foo", func(s string) (bool, string) { return len(s) == 3, "" })
			cmd.Options = Options{Each: true, Strict: true}

			So(match(cmd, "foo bar baz"), ShouldBeTrue)
			So(match(cmd, "foo foo true"), ShouldBeFalse) // strict
			So(match(cmd, "noye: test this out"), ShouldBeFalse)
		})

		Convey("match respond with mulitple parts", func() {
			cmd := newCommand("foo", func(s string) (bool, string) { return len(s) == 3, "" })
			cmd.Options = Options{Each: true, Respond: true}

			So(match(cmd, "noye: foo bar baz"), ShouldBeTrue)
			So(match(cmd, "noye: foo bar asdf"), ShouldBeTrue) // non strict
			So(match(cmd, "foo bar asdf"), ShouldBeFalse)
			So(match(cmd, "foo bar baz"), ShouldBeFalse)
			So(match(cmd, "noye: bar foo"), ShouldBeFalse)
		})

		Convey("match respond with mulitple parts, strict", func() {
			cmd := newCommand("foo", func(s string) (bool, string) { return len(s) == 3, "" })
			cmd.Options = Options{Each: true, Strict: true, Respond: true}

			So(match(cmd, "noye: foo bar baz"), ShouldBeTrue)
			So(match(cmd, "noye: foo bar asdf"), ShouldBeFalse) // strict
			So(match(cmd, "foo bar asdf"), ShouldBeFalse)
			So(match(cmd, "foo bar baz"), ShouldBeFalse)
			So(match(cmd, "noye: bar foo"), ShouldBeFalse)
		})

		Convey("match simple with a result", func() {
			cmd := newCommand("foo", func(s string) (bool, string) {
				if s == "test" {
					return true, "bar"
				}
				return false, ""
			})

			So(match(cmd, "foo test"), ShouldBeTrue)

			res := cmd.Results()
			So(len(res), ShouldEqual, 1)
			So(res[0], ShouldEqual, "bar")
		})

		Convey("match multiple with results", func() {
			cmd := newCommand("foo", func(s string) (bool, string) {
				ok, _ := regexp.MatchString("[0-9]", s)
				if ok {
					return ok, s
				}

				return false, ""
			})
			cmd.Options.Each = true

			So(match(cmd, "foo 1 0 0 4"), ShouldBeTrue)

			res := cmd.Results()
			So(len(res), ShouldEqual, 4)
			So(res, ShouldResemble, []string{"1", "0", "0", "4"})
		})

		Convey("match with built-in matchers", func() {
			Convey("using the simple matcher", func() {
				cmd := newCommand("foo", SimpleMatcher("bar"))
				cmd.Options.Each = true

				So(match(cmd, "foo bar bar bar"), ShouldBeTrue)

				res := cmd.Results()
				So(len(res), ShouldEqual, 0) // simplematcher doesn't capture
			})

			Convey("using the string matcher", func() {
				cmd := newCommand("foo", StringMatcher("bar", true))
				cmd.Options.Each = true

				So(match(cmd, "foo bar bar bar"), ShouldBeTrue)

				res := cmd.Results()
				So(len(res), ShouldEqual, 3)
				So(res, ShouldResemble, []string{"bar", "bar", "bar"})
			})

			Convey("using the regex matcher", func() {
				cmd := newCommand("foo", RegexMatcher(regexp.MustCompile("[0-9]"), true))
				cmd.Options.Each = true

				So(match(cmd, "foo a 1 2 b 3 c"), ShouldBeTrue)

				res := cmd.Results()
				So(len(res), ShouldEqual, 3)
				So(res, ShouldResemble, []string{"1", "2", "3"})
			})

			Convey("using the regex matcher, with no command", func() {
				re := regexp.MustCompile(`^(\w+:\/\/[\w@][\w.:@]+\/?[\w\.?=%&=\-@/$,]*)$`)
				cmd := newCommand("", RegexMatcher(re, true))
				cmd.Options.Each = true

				So(match(cmd, "http://google.com"), ShouldBeTrue)
				So(cmd.Results(), ShouldResemble, []string{"http://google.com"})
			})
		})
	})
}
