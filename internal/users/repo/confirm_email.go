package repo

import (
	"context"
	"fmt"

	"github.com/knstch/subtrack-libs/tracing"

	"users-service/internal/domain/enum"
)

func (r *DBRepo) ConfirmEmail(ctx context.Context, userID uint) error {
	ctx, span := tracing.StartSpan(ctx, "repo: ConfirmEmail")
	defer span.End()

	if err := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).
		Update("role", enum.VerifiedUserRole).Error; err != nil {
		return fmt.Errorf("db.Update: %w", err)
	}

	return nil
}
