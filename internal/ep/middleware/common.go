package middleware

import (
	"net/http"
	"regexp"
)

// little helper for http response info

type CommonResponseInfo struct {
	StatusCode int
	Size       int64
}

func EmptyCommonResponseInfo() *CommonResponseInfo {
	return NewCommonResponseInfo(-1, -1)
}

func NewCommonResponseInfo(statusCode int, size int64) *CommonResponseInfo {
	return &CommonResponseInfo{
		StatusCode: statusCode,
		Size:       size,
	}
}

// CommonResponseWriter http response logging
type CommonResponseWriter struct {
	http.ResponseWriter
	Info *CommonResponseInfo
}

func NewCommonResponseWriter(rw http.ResponseWriter) *CommonResponseWriter {
	return &CommonResponseWriter{
		ResponseWriter: rw,
		Info:           EmptyCommonResponseInfo(),
	}
}

// http.ResponseWriter =========================

func (clw *CommonResponseWriter) WriteHeader(statusCode int) {
	clw.ResponseWriter.WriteHeader(statusCode)
	clw.Info.StatusCode = statusCode
}

func (clw *CommonResponseWriter) Write(b []byte) (int, error) {
	size, err := clw.ResponseWriter.Write(b)
	clw.Info.Size += int64(size)

	return size, err
}

// =============================================

type PathMatcher struct {
	Method  string
	Path    string
	Pattern string
	matcher *regexp.Regexp
}

func NewPathMatcher(method string, path string, pattern string) *PathMatcher {
	matcher, err := regexp.Compile(path)
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

const PatternGetRoot string = `/[^\/]+`
const PatternPostRoot string = `/`
const Pattern
