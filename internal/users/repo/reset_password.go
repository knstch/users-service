package repo

import (
	"context"
	"fmt"
)

func (r *DBRepo) ResetPassword(ctx context.Context, email, password string) error {
	if err := r.db.Model(&User{}).WithContext(ctx).Where("email = ?", email).Update("password", password).Error; err != nil {
		return fmt.Errorf("db.Update: %w", err)
	}

	return nil
}
