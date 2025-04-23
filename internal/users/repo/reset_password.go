package repo

import (
	"context"
	"fmt"

	"github.com/knstch/subtrack-libs/tracing"
)

func (r *DBRepo) ResetPassword(ctx context.Context, email, password string) error {
	ctx, span := tracing.StartSpan(ctx, "repo: ResetPassword")
	defer span.End()

	if err := r.db.Model(&User{}).WithContext(ctx).Where("email = ?", email).Update("password", password).Error; err != nil {
		return fmt.Errorf("db.Update: %w", err)
	}

	return nil
}
