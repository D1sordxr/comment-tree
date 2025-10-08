package main

import (
	"context"
	"github.com/D1sordxr/comment-tree/internal/application/comment/usecase"
	"github.com/D1sordxr/comment-tree/internal/infrastructure/storage/postgres/repositories/comment/repo"
	"github.com/D1sordxr/comment-tree/internal/transport/http/api/comment/handler"
	"os"
	"os/signal"
	"syscall"

	loadApp "github.com/D1sordxr/comment-tree/internal/infrastructure/app"
	"github.com/D1sordxr/comment-tree/internal/infrastructure/config"
	"github.com/D1sordxr/comment-tree/internal/infrastructure/logger"
	"github.com/D1sordxr/comment-tree/internal/infrastructure/storage/postgres"
	"github.com/D1sordxr/comment-tree/internal/transport/http"

	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.NewApiConfig()

	log := logger.New(defaultLogger)
	log.Debug("Config data", "config", cfg)

	storageConn, err := dbpg.New(cfg.Storage.ConnectionString(), nil, nil)
	if err != nil {
		log.Error("Failed to connect to database", "error", err.Error())
		return
	}
	defer func() { _ = storageConn.Master.Close() }()
	if err = postgres.SetupStorage(storageConn.Master, cfg.Storage); err != nil {
		log.Error("Failed to setup storage", "error", err.Error())
		return
	}

	commentRepo := repo.New(storageConn)
	commentUC := usecase.New(log, commentRepo)
	commentHandler := handler.New(commentUC)

	httpServer := http.NewServer(
		log,
		&cfg.Server,
		commentHandler,
	)

	app := loadApp.NewApp(
		log,
		httpServer,
	)
	app.Run(ctx)
}

var defaultLogger zerolog.Logger

func init() {
	zlog.Init()
	defaultLogger = zlog.Logger
}
