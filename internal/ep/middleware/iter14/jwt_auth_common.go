package iter14

import (
	"fmt"

	"github.com/ElfAstAhe/url-shortener2/internal/ep/middleware"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
)

// AuthIter14PathMatcher describes watch paths
type AuthIter14PathMatcher struct {
	PathMatcher *middleware.PathMatcher
	StatusCode  int
}

func NewAuthIter14Path(method string, path string, pattern string, statusCode int) *AuthIter14PathMatcher {
	return &AuthIter14PathMatcher{
		PathMatcher: middleware.NewPathMatcher(method, path, pattern),
		StatusCode:  statusCode,
	}
}

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
			pathMatchers[matcher.PathMatcher.Method] = slice
		}

		slice = append(slice, matcher)
		pathMatchers[matcher.PathMatcher.Method] = slice
	}

	return &AuthIter14PathMatchers{
		PathMatchers: pathMatchers,
		log:          logger.GetLogger("authIter14PathMatchers"),
	}
}

func (pms *AuthIter14PathMatchers) HasPath(method string, path string) bool {
	return pms.GetPathMatcher(method, path) != nil
}

func (pms *AuthIter14PathMatchers) GetPathMatcher(method string, path string) *AuthIter14PathMatcher {
	slice, ok := pms.PathMatchers[method]
	if !ok {
		return nil
	}

	for _, item := range slice {
		if item.PathMatcher.Match(method, path) {
			// ToDo: remove logging
			pms.log.Info(fmt.Sprintf("Found matcher: method [%s] path [%s] pattern [%s] for request path [%s] ]", item.PathMatcher.Method, item.PathMatcher.Path, item.PathMatcher.Pattern, path))

			return item
		}
	}

	// ToDo: remove logging
	pms.log.Info(fmt.Sprintf("No matcher: for request path [%s] ]", path))

	return nil
}
