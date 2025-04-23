package users

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/knstch/subtrack-libs/svcerrs"
	"golang.org/x/crypto/bcrypt"

	"users-service/internal/users/repo"
)

func (svc *ServiceImpl) Login(ctx context.Context, email string, password string) (UserTokens, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return UserTokens{}, fmt.Errorf("mail.ParseAddress: %w", svcerrs.ErrInvalidData)
	}

	user, err := svc.repo.GetUser(ctx, repo.UserFilter{
		Email: email,
	})
	if err != nil {
		return UserTokens{}, fmt.Errorf("repo.GetUser: %w", err)
	}

	passwordFromDB, err := svc.repo.GetPassword(ctx, email)
	if err != nil {
		return UserTokens{}, fmt.Errorf("repo.GetPassword: %w", err)
	}

	if isMatch := svc.verifyPassword(passwordFromDB, password); !isMatch {
		return UserTokens{}, fmt.Errorf("wrong password: %w", svcerrs.ErrForbidden)
	}

	tokens, err := svc.mintJWT(user.UserID, user.Role)
	if err != nil {
		return UserTokens{}, fmt.Errorf("svc.mintJWT: %w", err)
	}

	return tokens, nil
}

func (svc *ServiceImpl) verifyPassword(hashedPassword, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+svc.passwordSecret)); err != nil {
		return false
	}

	return true
}
