package main

import (
	_ "NotesService/docs"
	"NotesService/internal/auth"
	"NotesService/internal/config"
	"NotesService/internal/handlers/note/deleteNote"
	"NotesService/internal/handlers/note/getAllNotes"
	"NotesService/internal/handlers/note/getOneNote"
	"NotesService/internal/handlers/note/putNote"
	"NotesService/internal/handlers/note/saveNotes"
	"NotesService/internal/handlers/users/registUser"
	"NotesService/internal/storage/postgresql"
	sl "NotesService/pkg/logger/logSlog"
	mwLogger "NotesService/pkg/logger/loggerMiddleware"
	logger "NotesService/pkg/logger/setupLogger"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/swaggo/http-swagger"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Notes Service API
// @version 1.0
// @description API for managing notes with JWT authentication
// @host localhost:8083
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	_ = godotenv.Load()

	cfg := config.MustLoad()
	fmt.Println(cfg)

	log := logger.SetupLogger(cfg.Env)
	log.Info("starting server", slog.String("env", cfg.Env))
	log.Debug("debug logging enabled")

	secret := os.Getenv("JWT_SECRET")
	// Время жизни токена
	tokenDuration := 30 * time.Minute

	jwtManager, err := auth.NewJWTManager(secret, tokenDuration)
	if err != nil {
		slog.Error("failed to create JWT manager", "error", err)
		os.Exit(1)
	}

	storage, err := postgresql.New(cfg.StoragePath())
	if err != nil {
		log.Error("error initializing storage", sl.Err(err))
		os.Exit(1)
	}

	//init router
	router := chi.NewRouter()

	//middleware
	router.Use(middleware.RequestID) //Генерирует уникальный ID для каждого запроса (для логов и отладки)
	router.Use(middleware.RealIP)    // Определяет реальный IP клиента (если есть прокси/балансировщик)
	router.Use(middleware.Logger)    //Логирует все запросы (URL, метод, статус)
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer) //Ловит паники (аварийные завершения) в хендлерах и не даёт упасть серверу
	router.Use(middleware.URLFormat) //Поддержка форматов URL вроде /api.json, /page.html

	// Редирект с /docs на /docs/index.html
	router.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/index.html", http.StatusMovedPermanently)
	})

	// Основной Swagger UI
	router.Get("/docs/*", httpSwagger.WrapHandler)

	router.Post("/users", registUser.New(log, storage, jwtManager))

	router.Route("/users/{id}/notes", func(r chi.Router) {
		r.Use(auth.JWTAuth(jwtManager))
		r.Post("/", saveNotes.New(log, storage))
		r.Get("/", getAllNotes.New(log, storage))
		r.Get("/{note_id}", getOneNote.New(log, storage))
		r.Put("/{note_id}", putNote.New(log, storage))
		r.Delete("/{note_id}", deleteNote.New(log, storage))
	})

	//START SERVER
	log.Info("starting server", slog.String("Address", cfg.HTTPServer.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("error starting server", sl.Err(err))
	}
	log.Error("stopping server", slog.String("Address", cfg.HTTPServer.Address))
}
