package middleware

import "net/http"

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
