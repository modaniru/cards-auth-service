package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/modaniru/cards-auth-service/internal/server"
	"github.com/modaniru/cards-auth-service/internal/service/auth"
	authservices "github.com/modaniru/cards-auth-service/internal/service/auth/auth_services"
	"github.com/modaniru/cards-auth-service/sqlc/db"
	"github.com/phsym/console-slog"
)

func main() {
	token := flag.String("t", "", "vk app token stub")
	flag.Parse()

	//TODO config file
	InitLogger("DEV")
	slog.Debug("logger init")

	conn, _ := sql.Open("postgres", "postgres://postgres:postgres@localhost:5555/postgres?sslmode=disable")
	db := db.New(conn)
	slog.Debug("database connect init")

	s := server.NewServer(
		&auth.AuthStub{},
		authservices.NewVKAuth(*token, db, conn),
	)
	slog.Debug("server init")
	slog.Debug("start server")
	http.ListenAndServe(":80", s.GetRouter())
}

// init logger [DEV, DEBUG, PROD]
func InitLogger(level string) {
	var handler slog.Handler
	switch level {
	case "PROD":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	case "DEBUG":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	default:
		handler = console.NewHandler(os.Stdout, &console.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
