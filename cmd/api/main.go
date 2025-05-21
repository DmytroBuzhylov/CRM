package main

import (
	"Test/config"
	taskHTTP "Test/internal/feature/task/delivery/http"
	"Test/internal/feature/task/repository/postgres"
	"Test/internal/feature/task/usecase"
	"Test/internal/pkg/db"
	"Test/internal/pkg/logger"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	logger.InnitLogger(cfg.Log.Level, cfg.Log.Production)
	log.Info().Msg("Application starting...")

	dbPool, err := db.InitPostgresSQLDB(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer dbPool.Close()

	taskRepo := postgres.NewTaskRepository(dbPool)
	createTaskUC := usecase.NewCreateTaskInteractor(taskRepo)
	getTaskUC := usecase.NewGetTaskInteractor(taskRepo)
	updateTaskUC := usecase.NewUpdateTaskInteractor(taskRepo)
	deleteTaskUC := usecase.NewDeleteTaskInteractor(taskRepo)

	taskHandler := taskHTTP.NewTaskHandler(createTaskUC, getTaskUC, updateTaskUC, deleteTaskUC)
	if cfg.Log.Production {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	//router.Use(taskHTTP.RequestLoggerMiddleware())
	router.Use(gin.Recovery())
	//router.Use(taskHTTP.ErrorHandlerMiddleware())

	apiV1 := router.Group("/api/v1")
	{
		tasksAPI := apiV1.Group("/tasks")
		{
			tasksAPI.POST("/", taskHandler.CreateTask)
			tasksAPI.GET("/:id", taskHandler.GetTask)
			tasksAPI.GET("/", taskHandler.GetTasks)
			tasksAPI.PUT("/:id", taskHandler.UpdateTask)
			tasksAPI.PATCH("/:id/status", taskHandler.UpdateTaskStatus)
			tasksAPI.PATCH("/:id/deadline", taskHandler.UpdateTaskDeadline)
			tasksAPI.DELETE("/:id", taskHandler.DeleteTask)
		}

	}

	router.GET("/healthz", func(c *gin.Context) {
		if err := dbPool.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "unhealthy", "database_error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
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

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}
	log.Info().Msg("Server exited gracefully")
}
