package main

import (
	"github.com/Vyacheslav1557/ms-auth/internal/app"
	"github.com/Vyacheslav1557/ms-auth/internal/config"
	"github.com/Vyacheslav1557/ms-auth/internal/services"
	"github.com/Vyacheslav1557/ms-auth/internal/storage"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func setupLog(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "dev":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		panic("")
	}

	return log
}

func main() {
	config.SetupCfg()
	cfg := config.Cfg()

	log := setupLog(cfg.Env)

	db, err := storage.New(cfg.DSN)
	if err != nil {
		panic(err)
	}
	authService := services.New(log, db)
	application := app.New(log, authService, cfg.Port)
	go func() {
		err = application.Run()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.Stop()
	log.Info("Gracefully stopped")
}
