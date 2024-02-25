package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/urfave/negroni"
)

func (s *Server) logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//wrap writer
		lrw := negroni.NewResponseWriter(w)
		//request start time
		start := time.Now()
		defer func() {
			requestTime := time.Since(start)
			status := lrw.Status()
			observe(requestTime, status)

			if status >= 400 {
				slog.Error(
					"user request",
					"url", r.RequestURI,
					"method", r.Method,
					"status", status,
					"time", requestTime.Microseconds(),
				)
				return
			}

			slog.Info(
				"user request",
				"url", r.RequestURI,
				"method", r.Method,
				"status", status,
				"time", requestTime.Microseconds(),
			)
		}()
		next.ServeHTTP(lrw, r)
		//request end time
	})
}

func (s *Server) jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")[7:]
		id, err := s.jwtService.ParseJwt(token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "id", id)) //TODO
		next.ServeHTTP(w, r)
	})
}
