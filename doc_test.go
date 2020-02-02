package access

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func ExampleMiddleware_stdlib() {
	logFunc := func(r *http.Request, status, size int, duration time.Duration) {
		log.Printf("[%s] %d %s %d %d", r.Method, status, r.URL.String(), size, duration.Milliseconds())
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hi!")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: Middleware(logFunc)(h),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panicf("web server: %v", err)
	}
}
