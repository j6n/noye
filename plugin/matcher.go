package plugin

import "regexp"

// Matcher is a function that takes a string
// returns whether it matched and the matched string
type Matcher func(string) (string, bool)

// NoopMatcher is a matcher that does nothing
func NoopMatcher() Matcher {
	return func(s string) (string, bool) { return s, true }
}

// SimpleMatcher is a matcher that matches the input
func SimpleMatcher(in string) Matcher {
	return StringMatcher(in, false)
}

// StringMatcher is a matcher that matches the input and also captures the input
func StringMatcher(in string, capture bool) Matcher {
	return func(s string) (res string, ok bool) {
		ok = s == in
		if capture && ok {
			res = s
		}

		return
	}
}

// RegexMatcher uses a regex to create a matcher
func RegexMatcher(re *regexp.Regexp, capture bool) Matcher {
	return func(s string) (res string, ok bool) {
		ok = re.MatchString(s)
		if capture && ok {
			res = s
		}

		return
	}
}
