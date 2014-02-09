package plugin

import "regexp"

type Matcher func(string) (bool, string)

func SimpleMatcher(in string) Matcher {
	return StringMatcher(in, false)
}

func StringMatcher(in string, capture bool) Matcher {
	return func(s string) (bool, string) {
		ok := s == in
		if capture && ok {
			return ok, s
		}

		return ok, ""
	}
}

func RegexMatcher(re *regexp.Regexp, capture bool) Matcher {
	return func(s string) (bool, string) {
		ok := re.MatchString(s)
		if capture && ok {
			return ok, s
		}

		return ok, ""
	}
}
