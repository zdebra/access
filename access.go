// Package access provides a simple middleware to implement access logging for
// a web server. It is accomplished by letting a caller implement its own
// logging function.
// 
// This is the example in cooperation with zap.Logger: 
//
// 	func loggingFunc(logger *zap.Logger) access.LogFunc {
// 		return func(r *http.Request, status, size int, duration time.Duration) {
// 			logger.Info("",
// 				zap.String("method", r.Method),
// 				zap.String("url", r.URL.String()),
// 				zap.Int("status", status),
// 				zap.Int("size", size),
// 				zap.Duration("duration", duration),
// 			)
// 		}
// 	}
//
// Then, this function is used to create logging middleware:
// 	access.Middleware(loggingFunc(logger))
// 
// This gives freedom to use user favorite logging solution.
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

// Middleware calls given LogFunc with request meta data
func Middleware(f LogFunc) func(next http.Handler) http.Handler {
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
