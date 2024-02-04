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
		slog.Info(
			"user request",
			"url", r.RequestURI,
			"method", r.Method, 
			"status", lrw.Status(),
			"time", end.Sub(start).Microseconds(), 
		)
	})
}