package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/knstch/subtrack-libs/svcerrs"
	"gorm.io/gorm"
)

func (r *DBRepo) UseRefreshToken(ctx context.Context, token string) (uint, error) {
	var tokens Token
	if err := r.db.Model(&Token{}).WithContext(ctx).Where("refresh_token = ? AND used = ?", token, false).First(&tokens).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("token not found: %w", svcerrs.ErrDataNotFound)
		}
		return 0, fmt.Errorf("db.First: %w", err)
	}

	if err := r.db.Model(&Token{}).WithContext(ctx).Where("id = ?", tokens.ID).Update("used", true).Error; err != nil {
		return 0, fmt.Errorf("db.Update: %w", err)
	}

	return tokens.UserID, nil
}
