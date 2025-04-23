package users

import (
	"context"
	"fmt"

	"users-service/internal/domain/dto"
	"users-service/internal/users/repo"

	"github.com/knstch/subtrack-libs/validator"
)

func (svc *ServiceImpl) GetUser(ctx context.Context, userID uint) (dto.User, error) {
	if err := validator.ValidateID(userID); err != nil {
		return dto.User{}, fmt.Errorf("validator.ValidateID: %w", err)
	}

	user, err := svc.repo.GetUser(ctx, repo.UserFilter{UserID: userID})
	if err != nil {
		return dto.User{}, fmt.Errorf("repo.GetUser: %w", err)
	}

	return dto.User{
		UserID: userID,
		Email:  user.Email,
		Role:   user.Role,
	}, nil
}
