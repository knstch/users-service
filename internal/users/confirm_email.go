package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/knstch/subtrack-libs/auth"
	"github.com/knstch/subtrack-libs/svcerrs"
	"github.com/knstch/subtrack-libs/tracing"

	"users-service/internal/domain/enum"
	"users-service/internal/users/repo"
	"users-service/internal/users/validator"
)

func (svc *ServiceImpl) ConfirmEmail(ctx context.Context, code string) (UserTokens, error) {
	ctx, span := tracing.StartSpan(ctx, "service: ConfirmEmail")
	defer span.End()

	if err := validator.ValidateConfirmationCode(code); err != nil {
		return UserTokens{}, fmt.Errorf("validator.ValidateConfirmationCode: %w", err)
	}

	userData, err := auth.GetUserData(ctx)
	if err != nil {
		return UserTokens{}, fmt.Errorf("auth.GetUserData: %w", err)
	}

	key := confirmationKey(userData.UserID)

	codeFromDB, err := svc.redis.Get(key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return UserTokens{}, fmt.Errorf("redis.Get: %w", svcerrs.ErrGone)
		}
		return UserTokens{}, fmt.Errorf("redis.Get: %w", err)
	}

	if codeFromDB != code {
		return UserTokens{}, fmt.Errorf("wrong confirmation code: %w", svcerrs.ErrInvalidData)
	}

	var userTokens UserTokens
	err = svc.repo.Transaction(func(st repo.Repository) error {
		if err = st.ConfirmEmail(ctx, userData.UserID); err != nil {
			return fmt.Errorf("st.ConfirmEmail: %w", err)
		}

		if err = st.DeactivateTokens(ctx, userData.UserID); err != nil {
			return fmt.Errorf("st.DeactivateTokens: %w", err)
		}

		userTokens, err = svc.mintJWT(userData.UserID, enum.VerifiedUserRole)
		if err != nil {
			return fmt.Errorf("svc.mintJWT: %w", err)
		}

		if err = st.StoreTokens(ctx, userData.UserID, userTokens.AccessToken, userTokens.RefreshToken); err != nil {
			return fmt.Errorf("st.StoreTokens: %w", err)
		}

		return nil
	})

	if err = svc.redis.Del(key).Err(); err != nil {
		return UserTokens{}, fmt.Errorf("redis.Del: %w", err)
	}

	return userTokens, nil
}
