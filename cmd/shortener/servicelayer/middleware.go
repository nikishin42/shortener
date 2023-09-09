package servicelayer

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
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

func (s *Server) Logging(h http.Handler) http.Handler {
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

func (s *Server) Decompressor(h http.Handler) http.Handler {
	decompressorFn := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		compressed := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
		if !compressed {
			h.ServeHTTP(w, r)
			return
		}
		zr, err := gzip.NewReader(r.Body)
		if err != nil {
			s.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		oldBody := r.Body
		defer oldBody.Close()
		r.Body = zr
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(decompressorFn)
}

type compressorWriter struct {
	http.ResponseWriter
	gz io.WriteCloser
}

func (c compressorWriter) Write(data []byte) (int, error) {
	return c.gz.Write(data)
}

func (s *Server) Compressor(h http.Handler) http.Handler {
	compressorFn := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.Logger.Error(err)
			return
		}
		defer gz.Close()
		h.ServeHTTP(compressorWriter{
			ResponseWriter: w,
			gz:             gz,
		}, r)
	}
	return http.HandlerFunc(compressorFn)
}
