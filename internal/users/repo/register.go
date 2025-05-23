package repo

import (
	"context"
	"fmt"

	"github.com/knstch/subtrack-libs/svcerrs"
	"github.com/knstch/subtrack-libs/tracing"
	"go.uber.org/zap"

	"users-service/internal/domain/enum"
)

func (r *DBRepo) CreateUser(ctx context.Context, email string, password string, role enum.Role) (uint, error) {
	ctx, span := tracing.StartSpan(ctx, "repo: CreateUser")
	defer span.End()

	user := &User{
		Email:    email,
		Password: password,
		Role:     role.String(),
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		r.lg.Error("error creating user", zap.Error(err), zap.String("method", "Register"))
		if isUniqueViolation(err) {
			return 0, fmt.Errorf("db.Create: %w", svcerrs.ErrConflict)
		}
		return 0, fmt.Errorf("db.Create: %w", err)
	}

	return user.ID, nil
}
