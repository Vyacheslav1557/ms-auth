package app

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"log/slog"
	"ms-auth/internal/config"
	"ms-auth/internal/database/postgresql"
	"ms-auth/internal/transport/rest/handlers/auth"
	"ms-auth/internal/transport/rest/handlers/users"
	"ms-auth/internal/transport/rest/middleware/logger"
	"net/http"
	"os"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// TODO: удаление неактивных в течение 60 дней пользователей
// TODO: исправить логи
// TODO: env, config, build
// TODO: написать client

func Run() {
	config.Init("./configs/config.yaml")

	cfg := config.Cfg()

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	log.Info("initializing server",
		slog.String("address", cfg.HTTPServer.Host+":"+cfg.HTTPServer.Port),
	)
	log.Debug("logger debug mode enabled")

	storage, err := postgresql.New(
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Storage.Host,
			cfg.Storage.Port,
			cfg.Storage.User,
			cfg.Storage.Password,
			cfg.Storage.DBname,
			cfg.Storage.SSLmode,
		),
	)
	if err != nil {
		panic(err.Error())
	}
	//--------------------------------------------R ROUTER--------------------------------------------------------------
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Set-Cookie"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(httprate.LimitByIP(100, 1*time.Minute))
	//-------------------------------------------AUTH ROUTER------------------------------------------------------------
	authRouter := chi.NewRouter()
	authRouter.Post("/login", auth.Login(log, storage))
	authRouter.Post("/logout", auth.Logout(log, storage))
	authRouter.Post("/refresh", auth.Refresh(log, storage))
	//------------------------------------------USERS ROUTER------------------------------------------------------------
	usersRouter := chi.NewRouter()
	usersRouter.Post("/", users.NewUser(log, storage))
	//-------------------------------------------MAIN ROUTER------------------------------------------------------------
	apiV1Router := chi.NewRouter()
	apiV1Router.Use(logger.NewLogger(log))
	apiV1Router.Mount("/auth", authRouter)
	apiV1Router.Mount("/users", usersRouter)
	//------------------------------------------------------------------------------------------------------------------
	r.Mount("/api/v1", apiV1Router)
	err = http.ListenAndServe(cfg.HTTPServer.Host+":"+cfg.HTTPServer.Port, r)
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
