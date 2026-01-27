package iter14

import (
	"strings"

	"github.com/ElfAstAhe/url-shortener2/internal/ep/middleware"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
)

// AuthIter14PathMatcher describes a watch path
type AuthIter14PathMatcher struct {
	PathMatcher           *middleware.PathMatcher
	InfoAbsentStatusCode  int
	InfoInvalidStatusCode int
	DefaultStatusCode     int
}

func NewAuthIter14PathMatcher(method string, path string, pattern string, infoAbsentStatusCode int, infoInvalidStatusCode int, defaultStatusCode int) *AuthIter14PathMatcher {
	return &AuthIter14PathMatcher{
		PathMatcher:           middleware.NewPathMatcher(method, path, pattern),
		InfoAbsentStatusCode:  infoAbsentStatusCode,
		InfoInvalidStatusCode: infoInvalidStatusCode,
		DefaultStatusCode:     defaultStatusCode,
	}
}

// AuthIter14PathMatchers describes watchable paths
type AuthIter14PathMatchers struct {
	PathMatchers map[string][]*AuthIter14PathMatcher
	log          logger.Logger
}

func NewAuthIter14PathMatchers(matchers []*AuthIter14PathMatcher, logger logger.Logger) *AuthIter14PathMatchers {
	pathMatchers := make(map[string][]*AuthIter14PathMatcher)

	for _, matcher := range matchers {
		slice, ok := pathMatchers[matcher.PathMatcher.Method]
		if !ok {
			slice = make([]*AuthIter14PathMatcher, 0)
		}

		if iter14AuthWatchPathExists(slice, matcher.PathMatcher.Method, matcher.PathMatcher.Path) {
			continue
		}

		pathMatchers[matcher.PathMatcher.Method] = append(slice, matcher)
	}

	return &AuthIter14PathMatchers{
		PathMatchers: pathMatchers,
		log:          logger.GetLogger("authIter14PathMatchers"),
	}
}

func (pms *AuthIter14PathMatchers) Match(method string, path string) bool {
	return pms.GetPathMatcher(method, path) != nil
}

func (pms *AuthIter14PathMatchers) GetPathMatcher(method string, path string) *AuthIter14PathMatcher {
	slice, ok := pms.PathMatchers[method]
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

func iter14AuthWatchPathExists(src []*AuthIter14PathMatcher, method string, path string) bool {
	if len(src) == 0 || strings.TrimSpace(method) == "" || strings.TrimSpace(path) == "" {
		return false
	}

	for _, item := range src {
		if item.PathMatcher.Method == method && item.PathMatcher.Path == path {
			return true
		}
	}

	return false
}
