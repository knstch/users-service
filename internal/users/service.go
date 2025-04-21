package users

import (
	"context"

	"go.uber.org/zap"

	"github.com/go-redis/redis"

	"users-service/config"
	"users-service/internal/users/repo"

	"github.com/knstch/subtrack-kafka/producer"
)

type ServiceImpl struct {
	lg *zap.Logger

	repo repo.Repository

	passwordSecret string
	jwtSecret      string

	redis    *redis.Client
	producer *producer.Producer
}

type Users interface {
	CreateUser(ctx context.Context, email string, password string) (UserTokens, error)
	ConfirmEmail(ctx context.Context, code string) (UserTokens, error)
	RefreshToken(ctx context.Context, refreshToken string) (UserTokens, error)
}

func NewService(lg *zap.Logger, repo repo.Repository, redis *redis.Client, cfg config.Config) *ServiceImpl {
	return &ServiceImpl{
		lg:             lg,
		repo:           repo,
		passwordSecret: cfg.PasswordSecret,
		jwtSecret:      cfg.JwtSecret,
		redis:          redis,
		producer:       producer.NewProducer(cfg.KafkaAddr),
	}
}
