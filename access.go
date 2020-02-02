package access

import (
	"net/http"
	"time"
)

type myResponseWriter struct {
	bytesWritten  int
	status        int
	headerWritten bool
	http.ResponseWriter
}

func (m *myResponseWriter) WriteHeader(code int) {
	if !m.headerWritten {
		m.status = code
		m.headerWritten = true
		m.ResponseWriter.WriteHeader(code)
	}
}

func (m *myResponseWriter) Write(buf []byte) (int, error) {
	m.WriteHeader(http.StatusOK)
	bw, err := m.ResponseWriter.Write(buf)
	m.bytesWritten += bw
	return bw, err
}

// LogFunc is a type of function which caller has to provide (use in logging func)
type LogFunc func(r *http.Request, status, size int, duration time.Duration)

// Handler calls given LogFunc with request meta data
func Handler(f LogFunc) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			mw := &myResponseWriter{
				ResponseWriter: w,
			}
			next.ServeHTTP(mw, r)
			f(r, mw.status, mw.bytesWritten, time.Since(start))
		})
	}
}
