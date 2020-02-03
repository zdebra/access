# Access

Access is a simple middleware to implement access logging for golang web server.

## Example usage

```go
import "go.uber.org/zap"

// loggingFunc implemented using zap.Logger
// you can use any logging system you want
func loggingFunc(logger *zap.Logger) access.LogFunc {
	return func(r *http.Request, status, size int, duration time.Duration) {
		logger.Info("",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.Int("status", status),
			zap.Int("size", size),
			zap.Duration("duration", duration),
		)
	}
}

func main() {
    logger, _ := zap.NewProduction()
    defer logger.Sync()
    
    h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "hi!")
    })

    srv := &http.Server{
        Addr:    ":8080",
        Handler: access.Middleware(loggingFunc(logger))(h),
    }

    if err := srv.ListenAndServe(); err != nil {
        logger.Panic("webserver stopped", zap.Error(err))
    }
}
```

## Alternatives

This project was inspired by [hlog](https://github.com/rs/zerolog/blob/master/hlog/hlog.go) which is
also an alternative. The problem with hlog is that it's utility package for zerolog and it has a 
dependency on it. It is also using another dependency just for wrapping `http.ResponseWriter` while
it can be easily achieved without it.
