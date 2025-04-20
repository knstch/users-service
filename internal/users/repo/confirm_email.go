package repo

import (
	"context"
	"fmt"

	"users-service/internal/domain/enum"
)

func (r *DBRepo) ConfirmEmail(ctx context.Context, userID uint) error {
	if err := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).
		Update("role", enum.VerifiedUserRole).Error; err != nil {
		return fmt.Errorf("db.Update: %w", err)
	}

	return nil
}
