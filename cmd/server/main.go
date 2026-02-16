package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/kerbatek/url-shortener/internal/handler"
	"github.com/kerbatek/url-shortener/internal/middleware"
	"github.com/kerbatek/url-shortener/internal/model"
	"github.com/kerbatek/url-shortener/internal/repository"
	"github.com/kerbatek/url-shortener/internal/service"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	log.Logger = logger

	var ctx = context.Background()
	var cfg model.Config
	var err error

	cfg.AppPort, err = strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil || cfg.AppPort == 0 {
		cfg.AppPort = 8080 // default port
	}
	cfg.DBName = os.Getenv("DB_NAME")
	cfg.DBUser = os.Getenv("DB_USER")
	cfg.DBPassword = os.Getenv("DB_PASSWORD")
	cfg.DBHost = os.Getenv("DB_HOST")
	cfg.DBPort, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil || cfg.DBPort == 0 {
		cfg.DBPort = 5432 // default PostgreSQL port
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		logger.Fatal().Err(err).Msg("Config parse failed")
	}

	config.MaxConns = 20
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Fatal().Err(err).Msg("Pool creation failed")
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Warn().Err(err).Msg("Database unreachable")
	}

	if err := runMigrations(ctx, pool, logger); err != nil {
		logger.Fatal().Err(err).Msg("Migration failed")
	}

	repo := repository.NewPostgresURLRepository(pool)
	svc := service.NewURLService(repo)
	h := handler.NewURLHandler(svc)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(middleware.Logger(logger))
	router.Use(gin.Recovery())
	router.StaticFile("/", "./static/index.html")
	router.Static("/static", "./static")
	router.POST("/shorten", h.ShortenURL)
	router.GET("/:code", h.RedirectURL)
	router.DELETE("/url/:id", h.DeleteURL)

	addr := fmt.Sprintf(":%d", cfg.AppPort)
	logger.Info().Str("addr", addr).Msg("Server starting")
	if err := router.Run(addr); err != nil {
		logger.Fatal().Err(err).Msg("Server failed")
	}
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool, logger zerolog.Logger) error {
	files, err := filepath.Glob("./migrations/*.up.sql")
	if err != nil {
		return fmt.Errorf("finding migrations: %w", err)
	}
	sort.Strings(files)

	for _, f := range files {
		sql, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("reading %s: %w", f, err)
		}
		if _, err := pool.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("executing %s: %w", f, err)
		}
		logger.Info().Str("file", f).Msg("Applied migration")
	}
	return nil
}
