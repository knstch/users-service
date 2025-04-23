package repo

import (
	"gorm.io/gorm"

	"users-service/internal/domain/enum"
)

type UserFilter struct {
	Email  string
	UserID uint
	Role   enum.Role
}

func (f *UserFilter) ToScope() func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Model(&User{})

		if f.Email != "" {
			tx = tx.Where("email = ?", f.Email)
		}

		if f.UserID != 0 {
			tx = tx.Where("id = ?", f.UserID)
		}

		if f.Role != "" {
			tx = tx.Where("role = ?", f.Role)
		}

		return tx
	}
}
