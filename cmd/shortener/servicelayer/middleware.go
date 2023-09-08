package servicelayer

import (
	"net/http"
	"time"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func (s *Server) WithLogging(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	h := http.HandlerFunc(f)

	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		respData := new(responseData)
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   respData,
		}
		h.ServeHTTP(&lw, r)
		duration := time.Since(start)

		s.Logger.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", respData.status,
			"duration", duration,
			"size", respData.size,
		)
	}
	return http.HandlerFunc(logFn)
}