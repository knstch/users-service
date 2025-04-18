package users

import (
	"context"
	"fmt"
	"math/rand"
	"net/mail"
	"time"

	"golang.org/x/crypto/bcrypt"

	"users-service/internal/users/repo"
	"users-service/internal/users/validator"

	"github.com/knstch/subtrack-libs/svcerrs"
)

type UserTokens struct {
	AccessToken  string
	RefreshToken string
}

func (svc *ServiceImpl) CreateUser(ctx context.Context, email string, password string) (UserTokens, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return UserTokens{}, fmt.Errorf("mail.ParseAddress: %w", svcerrs.ErrInvalidData)
	}

	if err = validator.ValidatePassword(ctx, password); err != nil {
		return UserTokens{}, fmt.Errorf("validator.ValidatePassword: %w", err)
	}

	passwordWithSalt, err := bcrypt.GenerateFromPassword([]byte(password+svc.passwordSecret), bcrypt.DefaultCost)
	if err != nil {
		return UserTokens{}, fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}

	userTokens := UserTokens{}
	if err = svc.repo.Transaction(func(st repo.Repository) error {
		userID, err := st.CreateUser(ctx, email, string(passwordWithSalt), "unverified_user")
		if err != nil {
			return fmt.Errorf("repo.CreateUser: %w", err)
		}

		userTokens.AccessToken, userTokens.RefreshToken, err = svc.mintJWT(userID, "unverified_user")
		if err != nil {
			return fmt.Errorf("svc.mintJWT: %w", err)
		}

		if err = st.StoreTokens(ctx, userID, userTokens.AccessToken, userTokens.RefreshToken); err != nil {
			return fmt.Errorf("st.StoreTokens: %w", err)
		}

		confirmationCode := rand.Int()
		if err = svc.redis.Set(fmt.Sprintf("confirmation-%d", userID), confirmationCode, time.Minute*30).Err(); err != nil {
			return fmt.Errorf("redis.Set: %w", err)
		}

		return nil
	}); err != nil {
		return UserTokens{}, fmt.Errorf("repo.Transaction: %w", err)
	}

	return userTokens, nil
}
