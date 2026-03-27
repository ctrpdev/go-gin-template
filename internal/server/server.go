package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api/internal/config"
	"api/internal/database"
	handler "api/internal/handler/http"
	"api/internal/logger"
	"api/internal/middleware"
	"api/internal/repository/postgres"
	"api/internal/repository/postgres/db"
	"api/internal/repository/redis"
	"api/internal/routes"
	"api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	redisClient "github.com/redis/go-redis/v9"
)

type Server struct {
	router *gin.Engine
	cfg    *config.Config
	dbPool *pgxpool.Pool
	rdb    *redisClient.Client
}

func New(cfg *config.Config) *Server {
	// Initialize Logger
	logger.Init(cfg.Environment)
	slog.Info("Initializing server components")

	// Database Connection
	dbPool, err := database.NewPostgresPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("Unable to connect to database", "err", err)
		os.Exit(1)
	}

	// Run Migrations automatically on startup
	database.RunDBMigrations("file://migrations", cfg.DatabaseURL)

	// Redis Connection
	rdb := database.NewRedisClient(cfg.RedisURL)

	// Dependency Injection
	queries := db.New(dbPool)
	userRepo := postgres.NewUserRepository(queries)
	sessionRepo := redis.NewSessionRepository(rdb)

	jwtSecret := []byte(cfg.JWTSecret)
	if len(jwtSecret) == 0 {
		slog.Warn("Using default JWT Secret. DO NOT use in production!")
		jwtSecret = []byte("my_super_secret_key")
	}

	userService := service.NewUserService(userRepo, sessionRepo, jwtSecret)
	userHandler := handler.NewUserHandler(userService)
	authMiddleware := middleware.NewAuthMiddleware(jwtSecret, sessionRepo)

	// Setup Router
	router := routes.SetupRouter(userHandler, authMiddleware)

	return &Server{
		router: router,
		cfg:    cfg,
		dbPool: dbPool,
		rdb:    rdb,
	}
}

func (s *Server) Run() {
	defer s.dbPool.Close()
	defer s.rdb.Close()

	srv := &http.Server{
		Addr:    s.cfg.ServerAddress,
		Handler: s.router,
	}

	go func() {
		slog.Info("Server listening", "address", s.cfg.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to run server", "err", err)
			os.Exit(1)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "err", err)
	}
	slog.Info("Server exiting")
}
