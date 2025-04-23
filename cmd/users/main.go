package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/go-redis/redis"
	"github.com/knstch/subtrack-libs/endpoints"
	"github.com/knstch/subtrack-libs/tracing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"users-service/config"
	"users-service/internal/endpoints/public"
	"users-service/internal/users"
	"users-service/internal/users/repo"

	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	if err := run(); err != nil {
		log.Println(err)
		recover()
	}
}

func run() error {
	args := os.Args

	dir, err := filepath.Abs(filepath.Dir(args[0]))
	if err != nil {
		return fmt.Errorf("filepath.Abs: %w", err)
	}

	if err := config.InitENV(dir); err != nil {
		return fmt.Errorf("config.InitENV: %w", err)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("config.GetConfig: %w", err)
	}

	shutdown := tracing.InitTracer(cfg.ServiceName, cfg.JaegerHost)
	defer shutdown(context.Background())

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(&lumberjack.Logger{
			Filename:   `./log/` + cfg.ServiceName + `_logfile.log`,
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
		}), zap.InfoLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(&lumberjack.Logger{
			Filename:   `./log/` + cfg.ServiceName + `_error.log`,
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
		}), zap.ErrorLevel),
	)
	lg := zap.New(core)

	dsnRedis, err := redis.ParseURL(cfg.GetRedisDSN())
	if err != nil {
		return err
	}
	redisClient := redis.NewClient(dsnRedis)

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("gorm.Open: %w", err)
	}
	dbRepo := repo.NewDBRepo(lg, db)

	svc := users.NewService(lg, dbRepo, redisClient, *cfg)

	publicController := public.NewController(svc, lg, cfg)
	publicEndpoints := endpoints.InitHttpEndpoints(cfg.ServiceName, publicController.Endpoints())

	srv := http.Server{
		Addr: ":" + cfg.PublicHTTPAddr,
		Handler: http.TimeoutHandler(
			publicEndpoints,
			time.Second*5,
			"service temporary unavailable",
		),
		ReadHeaderTimeout: time.Millisecond * 500,
		ReadTimeout:       time.Minute * 5,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err = srv.Shutdown(context.Background()); err != nil {
			log.Print(err)
		}
		close(idleConnsClosed)
	}()

	if err = srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Print(err)
	}

	<-idleConnsClosed

	return nil
}
