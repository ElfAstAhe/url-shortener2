package iter14

import (
	"github.com/ElfAstAhe/url-shortener2/internal/ep/middleware"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
)

// AuthIter14PathMatcher describes watch path
type AuthIter14PathMatcher struct {
	PathMatcher              *middleware.PathMatcher
	InfoAbsentStatusCode     int
	InfoInvalidStatusCode    int
	ErrAndUserInfoStatusCode int
	DefaultStatusCode        int
}

func NewAuthIter14PathMatcher(method string, path string, pattern string, infoAbsentStatusCode int, infoInvalidStatusCode int, errAndUserInfoStatusCode int, defaultStatusCode int) *AuthIter14PathMatcher {
	return &AuthIter14PathMatcher{
		PathMatcher:              middleware.NewPathMatcher(method, path, pattern),
		InfoAbsentStatusCode:     infoAbsentStatusCode,
		InfoInvalidStatusCode:    infoInvalidStatusCode,
		ErrAndUserInfoStatusCode: errAndUserInfoStatusCode,
		DefaultStatusCode:        defaultStatusCode,
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
			pathMatchers[matcher.PathMatcher.Method] = slice
		}

		if pathExists(slice, matcher.PathMatcher.Method, matcher.PathMatcher.Path) {
			continue
		}

		slice = append(slice, matcher)
		pathMatchers[matcher.PathMatcher.Method] = slice
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
			//			pms.log.Info(fmt.Sprintf("Found matcher: method [%s] path [%s] pattern [%s] for request path [%s] ]", item.PathMatcher.Method, item.PathMatcher.Path, item.PathMatcher.Pattern, path))

			return item
		}
	}

	//	pms.log.Info(fmt.Sprintf("No matcher: for request path [%s] ]", path))

	return nil
}

func pathExists(src []*AuthIter14PathMatcher, method string, path string) bool {
	if len(src) == 0 {
		return false
	}

	for _, item := range src {
		if item.PathMatcher.Method == method && item.PathMatcher.Path == path {
			return true
		}
	}

	return false
}
