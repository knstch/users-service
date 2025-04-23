package users

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/mail"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/knstch/subtrack-kafka/topics"
	"github.com/knstch/subtrack-libs/svcerrs"
	"github.com/knstch/users-api/event"
	"golang.org/x/crypto/bcrypt"

	"users-service/internal/users/repo"
	"users-service/internal/users/validator"
)

func resetPasswordKey(email string) string {
	return fmt.Sprintf("passwordReset-%s", email)
}

func (svc *ServiceImpl) ResetPassword(ctx context.Context, email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("mail.ParseAddress: %w", svcerrs.ErrInvalidData)
	}

	if _, err = svc.repo.GetUser(ctx, repo.UserFilter{
		Email: email,
	}); err != nil {
		return fmt.Errorf("repo.GetUser: %w", err)
	}

	confirmationCode := rand.Intn(9000) + 1000
	if err = svc.redis.Set(resetPasswordKey(email), confirmationCode, time.Minute*30).Err(); err != nil {
		return fmt.Errorf("redis.Set: %w", err)
	}

	payload, err := json.Marshal(&event.UserResetPassword{
		Email: email,
		Code:  strconv.Itoa(confirmationCode),
	})
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	if err = svc.repo.AddToOutbox(ctx, topics.TopicUserResetPassword, resetPasswordKey(email), payload); err != nil {
		return fmt.Errorf("repo.AddToOutbox: %w", err)
	}

	return nil
}

func (svc *ServiceImpl) ConfirmResetPassword(ctx context.Context, email, code, password string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("mail.ParseAddress: %w", svcerrs.ErrInvalidData)
	}

	if _, err = svc.repo.GetUser(ctx, repo.UserFilter{
		Email: email,
	}); err != nil {
		return fmt.Errorf("repo.GetUser: %w", err)
	}

	if err = validator.ValidateConfirmationCode(code); err != nil {
		return fmt.Errorf("validator.ValidateConfirmationCode: %w", err)
	}

	key := resetPasswordKey(email)

	codeFromDB, err := svc.redis.Get(key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("redis.Get: %w", svcerrs.ErrGone)
		}
		return fmt.Errorf("redis.Get: %w", err)
	}

	if codeFromDB != code {
		return fmt.Errorf("wrong confirmation code: %w", svcerrs.ErrInvalidData)
	}

	passwordWithSalt, err := bcrypt.GenerateFromPassword([]byte(password+svc.passwordSecret), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}

	if err = svc.repo.ResetPassword(ctx, email, string(passwordWithSalt)); err != nil {
		return fmt.Errorf("repo.ResetPassword: %w", err)
	}

	if err = svc.redis.Del(key).Err(); err != nil {
		return fmt.Errorf("redis.Del: %w", err)
	}

	return nil
}
