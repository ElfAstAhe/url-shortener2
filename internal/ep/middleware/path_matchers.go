package middleware

import "strings"

type PathMatchers struct {
	WatchPaths map[string][]*PathMatcher
}

func NewPathMatchers(matchers []*PathMatcher) *PathMatchers {
	watchPaths := make(map[string][]*PathMatcher)

	for _, matcher := range matchers {
		slice, ok := watchPaths[matcher.Method]
		if !ok {
			slice = make([]*PathMatcher, 0)
		}

		if watchPathExists(slice, matcher.Method, matcher.Path) {
			continue
		}

		watchPaths[matcher.Method] = append(slice, matcher)
	}

	return &PathMatchers{
		WatchPaths: watchPaths,
	}
}

func (pm *PathMatchers) Match(method string, path string) bool {
	return pm.GetPathMatcher(method, path) != nil
}

func (pm *PathMatchers) GetPathMatcher(method string, path string) *PathMatcher {
	slice, ok := pm.WatchPaths[method]
	if !ok {
		return nil
	}

	for _, item := range slice {
		if strings.TrimSpace(method) == item.Method && strings.TrimSpace(path) == item.Path {
			return item
		}
	}

	return nil
}

func watchPathExists(src []*PathMatcher, method string, path string) bool {
	if strings.TrimSpace(path) == "" || strings.TrimSpace(method) == "" || len(src) == 0 {
		return false
	}

	for _, item := range src {
		if item.Method == method && item.Path == path {
			return true
		}
	}

	return false
}
