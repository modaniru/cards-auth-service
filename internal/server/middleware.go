package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/urfave/negroni"
)

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//wrap writer
		lrw := negroni.NewResponseWriter(w)
		//request start time
		start := time.Now()
		next.ServeHTTP(lrw, r)
		//request end time
		end := time.Now()
		status := lrw.Status()
		if status >= 400 {
			slog.Error(
				"user request",
				"url", r.RequestURI,
				"method", r.Method,
				"status", status,
				"time", end.Sub(start).Microseconds(),
			)
			return
		}

		slog.Info(
			"user request",
			"url", r.RequestURI,
			"method", r.Method,
			"status", status,
			"time", end.Sub(start).Microseconds(),
		)
	})
}
