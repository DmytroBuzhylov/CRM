package app

import (
	"Test/config"
	"Test/internal/pkg/provider"
	"Test/pkg/db"
	"Test/pkg/logger"
	"Test/pkg/storage"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type App struct {
	Config       *config.AppConfig
	Logger       *zerolog.Logger
	DB           *pgxpool.Pool
	Redis        *redis.Client
	StorageS3    storage.Storage
	Router       *gin.Engine
	DataProvider provider.DataProvider
}

func NewApp(cfg *config.AppConfig) (*App, error) {
	logger.InnitLogger(cfg.Log.Level, cfg.Log.Production)

	dbPool, err := db.InitPostgresSQLDB(cfg.Database)
	if err != nil {
		return nil, errors.New("Failed to initialize database")
	}

	minioS3, err := storage.NewMinioStorage(cfg.Minio.Endpoint, cfg.Minio.AccessKeyId, cfg.Minio.SecretAccessKey, "bucket", cfg.Minio.SSL)
	if err != nil {
		return nil, errors.New("Failed to initialize minioS3")
	}

	redisClient, err := db.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		return nil, errors.New("Failed to initialize redis client")
	}

	dataProvider := provider.NewPostgresProvider(dbPool)

	r := gin.New()

	if cfg.Log.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	return &App{
		Config:       cfg,
		Logger:       nil,
		DB:           dbPool,
		Redis:        redisClient,
		Router:       r,
		StorageS3:    minioS3,
		DataProvider: dataProvider,
	}, nil
}

func (a *App) Run() {
	log.Info().Msg("Application starting...")
	a.setupDependencies()

	srv := &http.Server{
		Addr:         ":" + a.Config.Server.Port,
		Handler:      a.Router,
		ReadTimeout:  a.Config.Server.ReadTimeout,
		WriteTimeout: a.Config.Server.WriteTimeout,
		IdleTimeout:  a.Config.Server.IdleTimeout,
	}

	go func() {
		log.Info().Msgf("Server listening on %s", srv.Addr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed to listen and serve")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}
	log.Info().Msg("Server exited gracefully")
}
