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
	jwt2 "github.com/Khasmag06/kode-notes/pkg/jwt"
	logger2 "github.com/Khasmag06/kode-notes/pkg/logger"
	"github.com/Khasmag06/kode-notes/pkg/postgres"
	"github.com/Khasmag06/kode-notes/pkg/redis"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"net"
	"net/http"
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

	logger, err := logger2.New(cfg.Logger.LogFilePath, cfg.Logger.Level)
	if err != nil {
		log.Fatalf("failed to build logger: %s", err)
	}
	defer func() { _ = logger.Sync() }()

	ctx := context.Background()

	db, err := postgres.NewDB(ctx, cfg.PG)
	if err != nil {
		logger.Fatalf("failed to connect to postgres db: %s", err)
	}
	defer db.Close()

	redisDB, err := redis.ConnectRedis(ctx, cfg.Redis)
	if err != nil {
		logger.Fatalf("failed to connect to postgres redis db: %s", err)
	}

	jwt, err := jwt2.New(cfg.JWT.SignKey)
	if err != nil {
		logger.Fatal(err)
	}

	decoder, err := decoder2.New(cfg.Decoder.SecretKey)
	if err != nil {
		logger.Fatal(err)
	}

	hasher, err := hasher2.New(cfg.Hasher.Salt)
	if err != nil {
		logger.Fatal(err)
	}

	authRepo := authRepository.New(db.Pool)
	authService := auth.New(authRepo, hasher, jwt, decoder)

	noteRepo := noteRepository.New(db.Pool)
	noteCache := noteRepoWithCache.New(redisDB, noteRepo, logger)
	noteService := note.New(noteCache)
	yandexSpeller := speller.New(cfg.Speller.URL)

	r := api.NewHandler(authService, noteService, decoder, logger, yandexSpeller)

	server := http.Server{
		Addr:    net.JoinHostPort("", cfg.Server.Port),
		Handler: r,
	}
	logger.Info("Starting http server...")
	logger.Fatal(server.ListenAndServe())
}
