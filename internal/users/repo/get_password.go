package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/knstch/subtrack-libs/svcerrs"
	"github.com/knstch/subtrack-libs/tracing"
	"gorm.io/gorm"
)

func (r *DBRepo) GetPassword(ctx context.Context, email string) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "repo: GetPassword")
	defer span.End()

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
