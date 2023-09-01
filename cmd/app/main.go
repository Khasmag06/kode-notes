package main

import (
	"context"
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
	"net/http"

	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logger2.New(cfg.Logger.LogFilePath, cfg.Logger.Level)
	if err != nil {
		log.Fatalf("failed to build logger: %s", err)
	}
	defer func() { _ = logger.Sync() }()

	ctx := context.Background()

	db, err := postgres.NewDB(ctx, cfg.PG)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	redisDB, err := redis.ConnectRedis(ctx, cfg.Redis)
	if err != nil {
		log.Println(err)
	}

	jwt, err := jwt2.New(cfg.JWT.SignKey)
	if err != nil {
		logger.Fatal(err)
	}

	decoder, err := decoder2.New(cfg.Decoder.SecretKey)
	if err != nil {
		logger.Fatal(err)
	}
	hasher := hasher2.New(cfg.Hasher.Salt)

	authRepo := authRepository.New(db.Pool)
	authService := auth.New(authRepo, hasher, jwt, decoder)

	noteRepo := noteRepository.New(db.Pool)
	noteCache := noteRepoWithCache.New(redisDB, noteRepo, logger)
	noteService := note.New(noteCache)
	yandexSpeller := speller.New()

	r := api.NewHandler(authService, noteService, decoder, logger, yandexSpeller)

	logger.Fatal(http.ListenAndServe(cfg.Server.Port, r))
}
