package httpserver

import (
	"fmt"
	"movie-lib/pkg/logger"
	"net/http"
)

type ResponseWriterInterceptor struct {
	http.ResponseWriter
	StatusCode int
}

func (rw *ResponseWriterInterceptor) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func logMiddleware(next http.Handler, log logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &ResponseWriterInterceptor{w, http.StatusOK}
		next.ServeHTTP(rw, r)
		logStr := fmt.Sprintf("%s %s %d\n", r.Method, r.RequestURI, rw.StatusCode)
		if rw.StatusCode != http.StatusInternalServerError {
			log.InfoLog(logStr)
		} else {
			log.ErrorLog(logStr)
		}
	})
}
