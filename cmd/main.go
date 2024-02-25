package main

import (
	"database/sql"
	"github.com/modaniru/cards-auth-service/internal/storage"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	_ "net/http/pprof"
)

/*
todo migration
todo configuration
todo dockerfile
todo docker-compose
todo api gateway
*/

func main() {
	//token := flag.String("t", "", "vk app token stub")

	//TODO config file
	InitLogger("DEV")
	slog.Debug("logger init")

	token := os.Getenv("TOKEN")
	dataSource := os.Getenv("DATA_SOURCE")

	if token == "" {
		slog.Error("missing token")
		os.Exit(1)
	}
	slog.Debug("token was load")

	go http.ListenAndServe(":6060", nil)
	go func() {
		err := prometheus(":8082")
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	conn, _ := sql.Open("postgres", dataSource)
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

func prometheus(port string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(port, mux)
}
