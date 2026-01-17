package middleware

import "regexp"

type PathMatcher struct {
	Method  string
	Path    string
	Pattern string
	matcher *regexp.Regexp
}

func NewPathMatcher(method string, path string, pattern string) *PathMatcher {
	matcher, err := regexp.Compile(pattern)
	if err != nil {
		matcher = nil
	}

	return &PathMatcher{
		Method:  method,
		Path:    path,
		Pattern: pattern,
		matcher: matcher,
	}
}

func (m *PathMatcher) Match(method string, path string) bool {
	if m.matcher == nil {
		return false
	}

	if m.Method != method {
		return false
	}

	return m.matcher.MatchString(path)
}
