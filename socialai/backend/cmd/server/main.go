package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"socialai/internal/ai"
	"socialai/internal/api"
	"socialai/internal/auth"
	"socialai/internal/config"
	"socialai/internal/media"
	"socialai/internal/repository"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		logger.Error("load configuration", "error", err)
		os.Exit(1)
	}
	if os.Getenv("JWT_SECRET") == "" {
		logger.Warn("JWT_SECRET is unset; using an ephemeral development secret")
	}

	repo, err := buildRepository(cfg)
	if err != nil {
		logger.Error("initialize repository", "error", err)
		os.Exit(1)
	}
	defer repo.Close()

	storage, mediaHandler, err := buildStorage(cfg)
	if err != nil {
		logger.Error("initialize media storage", "error", err)
		os.Exit(1)
	}
	generator, err := buildGenerator(cfg)
	if err != nil {
		logger.Error("initialize AI image generator", "error", err)
		os.Exit(1)
	}

	server := api.New(repo, storage, generator, auth.NewTokenManager(cfg.JWTSecret, cfg.TokenTTL), api.Settings{
		AllowedOrigins: cfg.AllowedOrigins, MaxUploadBytes: cfg.MaxUploadBytes, TrustProxy: cfg.TrustProxy,
	}, logger)
	httpServer := &http.Server{
		Addr: cfg.Address, Handler: server.Handler(mediaHandler), ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout: 30 * time.Second, WriteTimeout: 2 * time.Minute, IdleTimeout: 60 * time.Second,
	}

	go func() {
		logger.Info("SocialAI API listening", "address", cfg.Address, "repository", cfg.RepositoryBackend, "media", cfg.MediaBackend, "ai", cfg.AIBackend)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("serve", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown", "error", err)
	}
}

func buildRepository(cfg config.Config) (repository.Repository, error) {
	switch cfg.RepositoryBackend {
	case "memory":
		return repository.NewMemory(), nil
	case "file":
		return repository.NewFile(cfg.DataFile)
	case "elasticsearch":
		return repository.NewElasticsearch(context.Background(), cfg.ElasticsearchURL, cfg.ElasticsearchUser, cfg.ElasticsearchPass, cfg.ElasticsearchPref)
	default:
		return nil, errors.New("REPOSITORY_BACKEND must be memory, file, or elasticsearch")
	}
}

func buildStorage(cfg config.Config) (media.Storage, http.Handler, error) {
	switch cfg.MediaBackend {
	case "local":
		local, err := media.NewLocal(cfg.MediaDirectory, cfg.PublicURL)
		if err != nil {
			return nil, nil, err
		}
		return local, local.Handler(), nil
	case "gcs":
		gcs, err := media.NewGCS(cfg.GCSBucket, cfg.GCSBearerToken)
		return gcs, nil, err
	default:
		return nil, nil, errors.New("MEDIA_BACKEND must be local or gcs")
	}
}

func buildGenerator(cfg config.Config) (ai.Generator, error) {
	switch cfg.AIBackend {
	case "local":
		return ai.Local{}, nil
	case "openai":
		return ai.NewOpenAI(cfg.OpenAIAPIKey, cfg.OpenAIModel, cfg.OpenAIBaseURL)
	default:
		return nil, errors.New("AI_IMAGE_BACKEND must be local or openai")
	}
}
