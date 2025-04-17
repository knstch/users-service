package repo

import (
	"context"
	"time"
)

type Repository interface {
	CreateUser(ctx context.Context, email string, password string, role string) (uint, error)
	StoreTokens(ctx context.Context, userID uint, accessToken string, refreshToken string) error

	Transaction(fn func(st Repository) error) error
}

type User struct {
	ID        uint
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
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
}

func (UsersData) TableName() string {
	return "users_data"
}

type Tokens struct {
	ID           uint
	UserID       uint
	AccessToken  string
	RefreshToken string
	Used         bool
}

func (Tokens) TableName() string {
	return "tokens"
}
