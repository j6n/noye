package plugin

import (
	"log"
	"regexp"
	"strings"

	"github.com/j6n/noye/noye"
)

type Command struct {
	Respond bool
	Command string
	Each    bool
	Matcher func(string) bool
}

func (c Command) Match(msg noye.Message) bool {
	// split text in parts so we can drop nick/cmd if needed
	parts := strings.Fields(msg.Text)

	// check to see if a nick was prefixed
	if c.Respond {
		nick := "noye" // TODO get current nick
		ok, err := regexp.MatchString(`(?:`+nick+`[:,]?\s*)`, parts[0])
		if err != nil || !ok {
			log.Println("expected nick to match")
			return false
		}
	}

	// less typing for later
	hasCommand := c.Command != ""

	// if we expect a nick prefix and a command, but only have 1 part
	if (len(parts) == 1 && c.Respond) && hasCommand {
		log.Println("have a command but expected more parts")
		return false
	}

	index := 0
	// skip first element if we've matched against `respond`
	if c.Respond {
		index++
	}

	// if we have a command, check to see if it matches the next part
	if hasCommand && !strings.EqualFold(c.Command, parts[index]) {
		log.Println("have a command but command and", index, parts[index], "doesn't match")
		return false
	}

	// skip next element if we've matched against `respond`
	if c.Respond && hasCommand {
		index++
	}

	// if no default matcher was provided, give them one that always returns true
	if c.Matcher == nil {
		log.Println("setting default matcher")
		c.Matcher = func(string) bool { return true }
	}

	// if we're using the matcher against each part
	if c.Each {
		// ...then match each remaining part
		for _, part := range parts[index:] {
			if !c.Matcher(part) {
				return false
			}
		}

		// nothing else to do, we've succeeded
		return true
	}

	// match against the parts rejoined as a string
	return c.Matcher(strings.Join(parts[index:], " "))
}
