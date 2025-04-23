package users

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/mail"
	"strconv"
	"time"

	"github.com/knstch/subtrack-kafka/topics"
	"github.com/knstch/subtrack-libs/tracing"
	"github.com/knstch/users-api/event"
	"golang.org/x/crypto/bcrypt"

	"users-service/internal/domain/enum"
	"users-service/internal/users/repo"
	"users-service/internal/users/validator"

	"github.com/knstch/subtrack-libs/svcerrs"
)

func confirmationKey(userID uint) string {
	return fmt.Sprintf("confirmation-%d", userID)
}

type UserTokens struct {
	AccessToken  string
	RefreshToken string
}

func (svc *ServiceImpl) Register(ctx context.Context, email string, password string) (UserTokens, error) {
	ctx, span := tracing.StartSpan(ctx, "service: Register")
	defer span.End()

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

	var (
		userTokens UserTokens
		userID     uint
	)

	if err = svc.repo.Transaction(func(st repo.Repository) error {
		userID, err = st.CreateUser(ctx, email, string(passwordWithSalt), enum.UnverifiedUserRole)
		if err != nil {
			return fmt.Errorf("repo.Register: %w", err)
		}

		userTokens, err = svc.mintJWT(userID, enum.UnverifiedUserRole)
		if err != nil {
			return fmt.Errorf("svc.mintJWT: %w", err)
		}

		if err = st.StoreTokens(ctx, userID, userTokens.AccessToken, userTokens.RefreshToken); err != nil {
			return fmt.Errorf("st.StoreTokens: %w", err)
		}

		return nil
	}); err != nil {
		return UserTokens{}, fmt.Errorf("repo.Transaction: %w", err)
	}

	confirmationCode := rand.Intn(9000) + 1000
	if err = svc.redis.Set(confirmationKey(userID), confirmationCode, time.Minute*30).Err(); err != nil {
		return UserTokens{}, fmt.Errorf("redis.Set: %w", err)
	}

	eventToOutbox, err := json.Marshal(&event.UserCreated{
		Email: email,
		Code:  strconv.Itoa(confirmationCode),
	})
	if err != nil {
		return UserTokens{}, fmt.Errorf("json.Marshal: %w", err)
	}

	if err = svc.repo.AddToOutbox(ctx, topics.TopicUserCreated, confirmationKey(userID), eventToOutbox); err != nil {
		return UserTokens{}, fmt.Errorf("repo.AddToOutbox: %w", err)
	}

	return userTokens, nil
}
