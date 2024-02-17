package main

import (
	"database/sql"
	"github.com/modaniru/cards-auth-service/internal/storage"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/modaniru/cards-auth-service/internal/server"
	"github.com/modaniru/cards-auth-service/internal/service/auth"
	authservices "github.com/modaniru/cards-auth-service/internal/service/auth/auth_services"
	jwtservice "github.com/modaniru/cards-auth-service/internal/service/jwt_service"
	"github.com/modaniru/cards-auth-service/sqlc/db"
	"github.com/phsym/console-slog"
)

func main() {
	//token := flag.String("t", "", "vk app token stub")

	//TODO config file
	InitLogger("DEV")
	slog.Debug("logger init")

	token := os.Getenv("TOKEN")
	if token == "" {
		slog.Error("missing token")
		os.Exit(1)
	}
	slog.Debug("token was load")

	conn, _ := sql.Open("postgres", "postgres://postgres:postgres@postgres/postgres?sslmode=disable")
	db := db.New(conn)
	slog.Debug("database connect init")

	globalStorage := storage.NewStorage(conn, db)
	slog.Debug("storage init")

	s := server.NewServer(
		jwtservice.NewJwtService("salt"),
		&auth.AuthStub{},
		authservices.NewVKAuth(token, globalStorage),
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
