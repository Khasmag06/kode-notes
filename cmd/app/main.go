package main

import (
	"context"
	"errors"
	"github.com/Khasmag06/kode-notes/config"
	"github.com/Khasmag06/kode-notes/internal/api"
	noteRepoWithCache "github.com/Khasmag06/kode-notes/internal/repository/note/cache"
	noteRepository "github.com/Khasmag06/kode-notes/internal/repository/note/postgres"
	authRepository "github.com/Khasmag06/kode-notes/internal/repository/user/postgres"
	"github.com/Khasmag06/kode-notes/internal/service/auth"
	"github.com/Khasmag06/kode-notes/internal/service/note"
	"github.com/Khasmag06/kode-notes/internal/service/speller"
	decoder2 "github.com/Khasmag06/kode-notes/pkg/decoder"
	hasher2 "github.com/Khasmag06/kode-notes/pkg/hasher"
	"github.com/Khasmag06/kode-notes/pkg/httpserver"
	jwt2 "github.com/Khasmag06/kode-notes/pkg/jwt"
	"github.com/Khasmag06/kode-notes/pkg/logger"
	"github.com/Khasmag06/kode-notes/pkg/postgres"
	"github.com/Khasmag06/kode-notes/pkg/redis"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @Title NoteService API
// @Description Сервис хранения заметок.
// @Version 1.0

// @BasePath /api
// @Host localhost:8080

// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name authorization

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.New("file://migrations", cfg.PG.URL)
	if err != nil {
		log.Fatalf("could not start sql migration... %v", err)
	}
	defer func() { _, _ = m.Close() }()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %v", err)
	}

	l, err := logger.New(cfg.Logger.LogFilePath, cfg.Logger.Level)
	if err != nil {
		log.Fatalf("failed to build logger: %s", err)
	}
	defer func() { _ = l.Sync() }()

	ctx := context.Background()

	db, err := postgres.NewDB(ctx, cfg.PG)
	if err != nil {
		l.Fatalf("failed to connect to postgres db: %s", err)
	}
	defer db.Close()

	redisDB, err := redis.ConnectRedis(ctx, cfg.Redis)
	if err != nil {
		l.Fatalf("failed to connect to postgres redis db: %s", err)
	}

	jwt, err := jwt2.New(cfg.JWT.SignKey)
	if err != nil {
		l.Fatal(err)
	}

	decoder, err := decoder2.New(cfg.Decoder.SecretKey)
	if err != nil {
		l.Fatal(err)
	}

	hasher, err := hasher2.New(cfg.Hasher.Salt)
	if err != nil {
		l.Fatal(err)
	}

	authRepo := authRepository.New(db.Pool)
	authService := auth.New(authRepo, hasher, jwt, decoder)

	noteRepo := noteRepository.New(db.Pool)
	noteCache := noteRepoWithCache.New(redisDB, noteRepo, l)
	noteService := note.New(noteCache)
	yandexSpeller := speller.New(cfg.Speller.URL)

	// HTTP Server
	l.Info("Starting http server...")
	handler := api.NewHandler(authService, noteService, decoder, l, yandexSpeller)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Errorf("app - Run - httpServer.Notify: %w", err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Errorf("app - Run - httpServer.Shutdown: %w", err)
	}

}
