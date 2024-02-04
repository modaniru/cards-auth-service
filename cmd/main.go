package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/modaniru/cards-auth-service/internal/server"
	"github.com/phsym/console-slog"
)

func main(){
	//TODO config file
	InitLogger("DEV")
	slog.Debug("logger init")
	s := server.NewServer()
	slog.Debug("server init")
	slog.Debug("start server")
	http.ListenAndServe(":80", s.GetRouter())
}

//init logger [DEV, DEBUG, PROD]
func InitLogger(level string){
	var handler slog.Handler
	switch level{
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