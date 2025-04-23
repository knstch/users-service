package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/knstch/subtrack-libs/svcerrs"
	"gorm.io/gorm"
)

func (r *DBRepo) GetPassword(ctx context.Context, email string) (string, error) {
	var user User
	if err := r.db.Model(&User{}).WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("user not found: %w", svcerrs.ErrDataNotFound)
		}
		return "", fmt.Errorf("db.First: %w", err)
	}
	r.db.Scopes()
	return user.Password, nil
}
