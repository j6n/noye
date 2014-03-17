package plugin

import "regexp"

// MatchFn is a function that takes a string and
// returns whether it matched and the matched string
type MatchFn func(string) (string, bool)

// Matcher returns a MatchFn which can match against different types of input
type Matcher interface {
	Match() MatchFn
}

// BaseMatcher is a matcher which implements Match() naively
type BaseMatcher struct{ Fn MatchFn }

// Match returns the MatchFn for the Matcher
func (b BaseMatcher) Match() MatchFn { return b.Fn }

// NoopMatcher is a matcher that does nothing
var NoopMatcher BaseMatcher = BaseMatcher{
	func(s string) (string, bool) { return s, true },
}

// SimpleMatcher is a matcher that matches the input to 'in'
type SimpleMatcher struct{ BaseMatcher }

// StringMatcher is a matcher which also captures the input
type StringMatcher struct{ BaseMatcher }

// RegexMatcher uses a regex to create a matcher
type RegexMatcher struct{ BaseMatcher }

// SimpleMatch returns a new SimpleMatcher
func SimpleMatch(in string) SimpleMatcher {
	return SimpleMatcher{BaseMatcher{StringMatch(in, false).Fn}}
}

// StringMatch returns a new StringMatcher
func StringMatch(in string, capture bool) StringMatcher {
	return StringMatcher{BaseMatcher{func(s string) (res string, ok bool) {
		ok = s == in
		if capture && ok {
			res = s
		}

		return
	}}}
}

// RegexMatch returns a new RegexMatcher
func RegexMatch(re *regexp.Regexp, capture bool) RegexMatcher {
	return RegexMatcher{BaseMatcher{func(s string) (res string, ok bool) {
		ok = re.MatchString(s)
		if capture && ok {
			res = s
		}

		return
	}}}
}
