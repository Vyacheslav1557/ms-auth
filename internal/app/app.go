package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log/slog"
	"ms-auth/internal/config"
	"ms-auth/internal/transport/rest/middleware/logger"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Run() {
	cfg := config.MustLoad("./configs/config.yaml")

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	log.Info("initializing server",
		slog.String("address", cfg.HTTPServer.Host+":"+cfg.HTTPServer.Port),
	)
	log.Debug("logger debug mode enabled")

	//storage, err := postgres.New(
	//	fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
	//		cfg.Storage.Host,
	//		cfg.Storage.Port,
	//		cfg.Storage.User,
	//		cfg.Storage.Password,
	//		cfg.Storage.DBname,
	//		cfg.Storage.SSLmode,
	//	),
	//)
	//if err != nil {
	//	panic(err.Error())
	//}
	//--------------------------------------------R ROUTER--------------------------------------------------------------
	r := chi.NewRouter()
	//-------------------------------------------MAIN ROUTER------------------------------------------------------------
	apiV1Router := chi.NewRouter()
	apiV1Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Set-Cookie"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))
	apiV1Router.Use(middleware.RequestID)
	apiV1Router.Use(middleware.Logger)
	apiV1Router.Use(middleware.Recoverer)
	apiV1Router.Use(middleware.URLFormat)
	apiV1Router.Use(logger.NewLogger(log))
	//------------------------------------------------------------------------------------------------------------------
	r.Mount("/api/v1", apiV1Router)

	err := http.ListenAndServe(cfg.HTTPServer.Host+":"+cfg.HTTPServer.Port, r)
	if err != nil {
		log.Error(err.Error())
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
