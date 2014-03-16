package plugin

import "regexp"

// Matcher is a function that takes a string and
// returns whether it matched and the matched string
type Matcher func(string) (bool, string)

// NoopMatcher is a matcher that does nothing
func NoopMatcher() Matcher {
	return func(s string) (bool, string) { return true, s }
}

// SimpleMatcher is a matcher that matches the input to 'in'
func SimpleMatcher(in string) Matcher {
	return StringMatcher(in, false)
}

// StringMatcher is a matcher which also captures the input
func StringMatcher(in string, capture bool) Matcher {
	return func(s string) (bool, string) {
		ok := s == in
		if capture && ok {
			return ok, s
		}

		return ok, ""
	}
}

// RegexMatcher uses a regex to create a matcher
func RegexMatcher(re *regexp.Regexp, capture bool) Matcher {
	return func(s string) (bool, string) {
		ok := re.MatchString(s)
		if capture && ok {
			return ok, s
		}

		return ok, ""
	}
}
