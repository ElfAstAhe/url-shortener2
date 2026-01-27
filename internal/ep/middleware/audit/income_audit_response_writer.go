package audit

import (
	"net/http"

	"github.com/ElfAstAhe/url-shortener2/internal/ep/middleware"
)

type IncomeAuditResponseWriter struct {
	http.ResponseWriter
	Info *middleware.CommonResponseInfo
	Data []byte
}

func NewIncomeAuditResponseWriter(rw http.ResponseWriter) *IncomeAuditResponseWriter {
	return &IncomeAuditResponseWriter{
		ResponseWriter: rw,
		Info:           middleware.EmptyCommonResponseInfo(),
		Data:           make([]byte, 0),
	}
}

// http.ResponseWriter =========================

func (iarw *IncomeAuditResponseWriter) WriteHeader(statusCode int) {
	iarw.ResponseWriter.WriteHeader(statusCode)
	iarw.Info.StatusCode = statusCode
}

func (iarw *IncomeAuditResponseWriter) Write(b []byte) (int, error) {
	size, err := iarw.ResponseWriter.Write(b)
	iarw.Info.Size += int64(size)
	iarw.Data = append(iarw.Data, b...)

	return size, err
}

// =============================================
