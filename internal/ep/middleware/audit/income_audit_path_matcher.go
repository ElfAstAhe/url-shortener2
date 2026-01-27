package audit

import (
	"github.com/ElfAstAhe/url-shortener2/internal/ep/middleware"
)

type IncomeAuditPathMatcher struct {
	PathMatcher      *middleware.PathMatcher
	urlExtractorFunc OriginalURLExtractorFunc
}

func NewIncomeAuditPathMatcher(method string, path string, pattern string, urlExtractorFunc OriginalURLExtractorFunc) *IncomeAuditPathMatcher {
	return &IncomeAuditPathMatcher{
		PathMatcher:      middleware.NewPathMatcher(method, path, pattern),
		urlExtractorFunc: urlExtractorFunc,
	}
}

type IncomeAuditPathMatchers struct {
	pathMatchers map[string][]*IncomeAuditPathMatcher
}

func NewIncomeAuditPathMatchers(matchers []*IncomeAuditPathMatcher) *IncomeAuditPathMatchers {
	pathMatchers := make(map[string][]*IncomeAuditPathMatcher)

	for _, matcher := range matchers {
		slice, ok := pathMatchers[matcher.PathMatcher.Method]
		if !ok {
			slice = make([]*IncomeAuditPathMatcher, 0)
		}

		if incomeAuditMatcherExists(slice, matcher.PathMatcher.Method, matcher.PathMatcher.Path) {
			continue
		}

		pathMatchers[matcher.PathMatcher.Method] = append(slice, matcher)
	}

	return &IncomeAuditPathMatchers{
		pathMatchers: pathMatchers,
	}
}

func (iapm *IncomeAuditPathMatchers) Match(method string, path string) bool {
	return iapm.GetPathMatcher(method, path) != nil
}

func (iapm *IncomeAuditPathMatchers) GetPathMatcher(method string, path string) *IncomeAuditPathMatcher {
	slice, ok := iapm.pathMatchers[method]
	if !ok {
		return nil
	}

	for _, item := range slice {
		if item.PathMatcher.Match(method, path) {
			return item
		}
	}

	return nil
}

func incomeAuditMatcherExists(pathMatchers []*IncomeAuditPathMatcher, method string, path string) bool {
	for _, matcher := range pathMatchers {
		if matcher.PathMatcher.Method == method && matcher.PathMatcher.Path == path {
			return true
		}
	}

	return false
}
