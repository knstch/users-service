package repo

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

func (r *DBRepo) StoreTokens(ctx context.Context, userID uint, accessToken string, refreshToken string) error {
	if err := r.db.WithContext(ctx).Create(&Tokens{
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Used:         false,
	}).Error; err != nil {
		r.lg.Error("error creating tokens", zap.Error(err), zap.String("method", "StoreTokens"))
		return fmt.Errorf("db.Create: %w", err)
	}

	return nil
}
