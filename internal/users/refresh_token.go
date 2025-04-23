package users

import (
	"context"
	"fmt"

	"github.com/knstch/subtrack-libs/tracing"

	"users-service/internal/users/repo"
)

func (svc *ServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (UserTokens, error) {
	ctx, span := tracing.StartSpan(ctx, "service: RefreshToken")
	defer span.End()

	var userTokens UserTokens
	err := svc.repo.Transaction(func(st repo.Repository) error {
		userID, err := st.UseRefreshToken(ctx, refreshToken)
		if err != nil {
			return fmt.Errorf("st.UseRefreshToken: %w", err)
		}

		user, err := st.GetUser(ctx, repo.UserFilter{
			UserID: userID,
		})
		if err != nil {
			return fmt.Errorf("st.GetUser: %w", err)
		}

		userTokens, err = svc.mintJWT(userID, user.Role)
		if err != nil {
			return fmt.Errorf("svc.mintJWT: %w", err)
		}

		if err = st.StoreTokens(ctx, userID, userTokens.AccessToken, userTokens.RefreshToken); err != nil {
			return fmt.Errorf("st.StoreTokens: %w", err)
		}

		return nil
	})
	if err != nil {
		return UserTokens{}, fmt.Errorf("repo.Transaction: %w", err)
	}

	return userTokens, nil
}
