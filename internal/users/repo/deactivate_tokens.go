package repo

import (
	"context"
	"fmt"
)

func (r *DBRepo) DeactivateTokens(ctx context.Context, userID uint) error {
	if err := r.db.WithContext(ctx).Model(&Token{}).Where("user_id = ?", userID).Update("used", true).Error; err != nil {
		return fmt.Errorf("db.Update: %w", err)
	}

	return nil
}
