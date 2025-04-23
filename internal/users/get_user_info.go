package users

import (
	"context"
	"fmt"

	"users-service/internal/domain/dto"
	"users-service/internal/users/repo"
)

func (svc *ServiceImpl) GetUserInfo(ctx context.Context, userID uint) (dto.User, error) {
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
