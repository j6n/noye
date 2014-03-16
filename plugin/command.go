package plugin

import (
	"regexp"
	"strings"

	"github.com/j6n/noye/noye"
)

// Command is a type that makes up a DSL.
// This DSL allows an irc-bot command, via chat to be matched
// in a simple, programmatic fashion.
type Command struct {
	Command  string
	Options  Options
	Matchers []Matcher

	results []string
}

type Options struct {
	Respond bool
	Strict  bool
	Each    bool
}

// Hear is a command that isn't directed toward the bot
// It takes a command string and a matcher and returns a Command
func Hear(cmd string, opt Options, matchers ...Matcher) *Command {
	return &Command{Command: cmd, Matchers: matchers}
}

// Respond is a command that is directed toward the bot
// It takes a command string and a matcher and returns a Command
func Respond(cmd string, opt Options, matchers ...Matcher) *Command {
	opt.Respond = true
	return &Command{Command: cmd, Options: opt, Matchers: matchers}
}

// Match matches the command to the noye.Message
// returning whether it matched or not
func (c *Command) Match(msg noye.Message) bool {
	// reset the results
	c.results = make([]string, 0)

	// split text in parts so we can drop nick/cmd if needed
	parts := strings.Fields(msg.Text)

	// check to see if a nick was prefixed
	if c.Options.Respond {
		nick := "noye" // TODO get current nick
		ok, err := regexp.MatchString(`(?:`+nick+`[:,]?\s*)`, parts[0])
		if err != nil || !ok {
			return false
		}
	}

	// less typing for later
	hasCommand := c.Command != ""

	// if we expect a nick prefix and a command, but only have 1 part
	if (len(parts) == 1 && c.Options.Respond) && hasCommand {
		return false
	}

	index := 0
	// skip first element if we've matched against `respond`
	if c.Options.Respond {
		index++
	}

	// if we have a command, check to see if it matches the next part
	if hasCommand && !strings.EqualFold(c.Command, parts[index]) {
		return false
	}

	// skip next element if we've matched against `respond`
	if c.Options.Respond || hasCommand {
		index++
	}

	// if no default matcher was provided, give them one that always returns true
	if len(c.Matchers) == 0 {
		c.Matchers = append(c.Matchers, func(string) (bool, string) { return true, "" })
	}

	// if we're using the matcher against each part
	if c.Options.Each {
		// ...then match each remaining part
		for _, part := range parts[index:] {
			// ...to each matcher
			for _, matcher := range c.Matchers {
				if ok, s := matcher(part); ok {
					if s != "" {
						c.results = append(c.results, s)
					}
				} else if c.Options.Strict {
					// if we are strict matching, then we've failed if we've found no match
					return false
				}
			}
		}

		// nothing else to do, we've succeeded
		return true
	}

	// match against the parts rejoined as a string
	input := strings.Join(parts[index:], " ")
	for _, matcher := range c.Matchers {
		if ok, s := matcher(input); ok {
			if s != "" {
				c.results = append(c.results, s)
			}
		} else if c.Options.Strict {
			return false
		}
	}

	return true
}

// Results returns the command results
func (c *Command) Results() []string {
	return c.results
}
