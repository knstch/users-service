package users

import (
	"context"

	"go.uber.org/zap"

	"github.com/go-redis/redis"

	"users-service/config"
	"users-service/internal/users/repo"
)

type ServiceImpl struct {
	lg *zap.Logger

	repo repo.Repository

	passwordSecret string
	jwtSecret      string

	redis *redis.Client
}

type Users interface {
	CreateUser(ctx context.Context, email string, password string) (UserTokens, error)
}

func NewService(lg *zap.Logger, repo repo.Repository, redis *redis.Client, cfg config.Config) *ServiceImpl {
	return &ServiceImpl{
		lg:             lg,
		repo:           repo,
		passwordSecret: cfg.PasswordSecret,
		jwtSecret:      cfg.JwtSecret,
		redis:          redis,
	}
}
