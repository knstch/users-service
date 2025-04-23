package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/knstch/subtrack-libs/svcerrs"
	"gorm.io/gorm"

	"users-service/internal/domain/dto"
	"users-service/internal/domain/enum"
)

func (r *DBRepo) GetUser(ctx context.Context, filter UserFilter) (dto.User, error) {
	var user User
	if err := r.db.WithContext(ctx).Scopes(filter.ToScope()).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.User{}, fmt.Errorf("user not found: %w", svcerrs.ErrDataNotFound)
		}
		return dto.User{}, fmt.Errorf("db.First: %w", err)
	}

	return dto.User{
		UserID: user.ID,
		Email:  user.Email,
		Role:   enum.Role(user.Role),
	}, nil
}
