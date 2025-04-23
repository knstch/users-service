package repo

import (
	"context"
	"fmt"

	"github.com/knstch/subtrack-libs/tracing"
)

func (r *DBRepo) DeactivateTokens(ctx context.Context, userID uint) error {
	ctx, span := tracing.StartSpan(ctx, "repo: DeactivateTokens")
	defer span.End()

	if err := r.db.WithContext(ctx).Model(&Token{}).Where("user_id = ?", userID).Update("used", true).Error; err != nil {
		return fmt.Errorf("db.Update: %w", err)
	}

	return nil
}
