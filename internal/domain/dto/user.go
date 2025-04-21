package dto

import "users-service/internal/domain/enum"

type User struct {
	UserID uint
	Email  string
	Role   enum.Role
}
