package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/knstch/subtrack-kafka/outbox"
	"github.com/knstch/subtrack-kafka/topics"

	"users-service/internal/domain/dto"
	"users-service/internal/domain/enum"
)

type Repository interface {
	CreateUser(ctx context.Context, email string, password string, role enum.Role) (uint, error)
	StoreTokens(ctx context.Context, userID uint, accessToken string, refreshToken string) error
	ConfirmEmail(ctx context.Context, userID uint) error
	DeactivateTokens(ctx context.Context, userID uint) error
	UseRefreshToken(ctx context.Context, token string) (uint, error)
	GetUser(ctx context.Context, filter UserFilter) (dto.User, error)
	GetPassword(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, email, password string) error

	Transaction(fn func(st Repository) error) error

	AddToOutbox(ctx context.Context, topic topics.KafkaTopic, key string, payload []byte) error
}

type User struct {
	ID        uint
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (User) TableName() string {
	return "users"
}

type UsersData struct {
	ID        uint
	UserID    uint
	Name      string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (UsersData) TableName() string {
	return "users_data"
}

type Token struct {
	ID           uint
	UserID       uint
	AccessToken  string
	RefreshToken string
	Used         bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

func (Token) TableName() string {
	return "tokens"
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func (r *DBRepo) AddToOutbox(ctx context.Context, topic topics.KafkaTopic, key string, payload []byte) error {
	if err := r.db.WithContext(ctx).Model(&outbox.Outbox{}).Create(&outbox.Outbox{
		Topic:   topic.String(),
		Key:     key,
		Payload: payload,
	}).Error; err != nil {
		return fmt.Errorf("db.Create: %w", err)
	}

	return nil
}
