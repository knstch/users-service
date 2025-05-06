package users

import (
	"context"
	"github.com/knstch/subtrack-libs/log"

	"github.com/go-redis/redis"

	"users-service/config"
	"users-service/internal/domain/dto"
	"users-service/internal/users/repo"

	"github.com/knstch/subtrack-kafka/producer"
)

type ServiceImpl struct {
	lg *log.Logger

	repo repo.Repository

	passwordSecret string
	jwtSecret      string

	redis    *redis.Client
	producer *producer.Producer
}

type Users interface {
	Register(ctx context.Context, email string, password string) (UserTokens, error)
	ConfirmEmail(ctx context.Context, code string) (UserTokens, error)
	RefreshToken(ctx context.Context, refreshToken string) (UserTokens, error)
	Login(ctx context.Context, email string, password string) (UserTokens, error)
	GetUserInfo(ctx context.Context, userID uint) (dto.User, error)
	ResetPassword(ctx context.Context, email string) error
	ConfirmResetPassword(ctx context.Context, email, code, password string) error
}

func NewService(lg *log.Logger, repo repo.Repository, redis *redis.Client, cfg config.Config) *ServiceImpl {
	return &ServiceImpl{
		lg:             lg,
		repo:           repo,
		passwordSecret: cfg.PasswordSecret,
		jwtSecret:      cfg.JwtSecret,
		redis:          redis,
		producer:       producer.NewProducer(cfg.KafkaAddr),
	}
}
