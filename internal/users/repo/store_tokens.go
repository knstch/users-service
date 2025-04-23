package repo

import (
	"context"
	"fmt"

	"github.com/knstch/subtrack-libs/tracing"
	"go.uber.org/zap"
)

func (r *DBRepo) StoreTokens(ctx context.Context, userID uint, accessToken string, refreshToken string) error {
	ctx, span := tracing.StartSpan(ctx, "repo: StoreTokens")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(&Token{
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
