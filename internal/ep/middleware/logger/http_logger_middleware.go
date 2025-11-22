package logger

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ElfAstAhe/url-shortener2/internal/ep/middleware"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
)

type HTTPReqRespLogger struct {
	log logger.Logger
}

func NewHTTPReqRespLogger(log logger.Logger) *HTTPReqRespLogger {
	return &HTTPReqRespLogger{
		log: log.GetLogger("logger middleware"),
	}
}

func (rrl *HTTPReqRespLogger) LogReqRes(nextHandler http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		crw := middleware.NewCommonResponseWriter(rw)

		nextHandler.ServeHTTP(crw, r)

		duration := time.Since(start)

		rrl.log.Infof("uri [%s] method [%s] duration [%v]ms status [%v] size [%v]",
			r.RequestURI, r.Method, strconv.FormatInt(duration.Milliseconds(), 10), crw.Info.StatusCode, crw.Info.Size)
	}

	return http.HandlerFunc(fn)
}
