package public

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/knstch/subtrack-libs/auth"
	"github.com/knstch/subtrack-libs/tracing"
	public "github.com/knstch/users-api/public"

	"users-service/internal/domain/enum"
)

func MakeGetUserInfoEndpoint(c *Controller) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return c.GetUserInfo(ctx, nil)
	}
}

func (c *Controller) GetUserInfo(ctx context.Context, _ *public.GetUserInfoRequest) (*public.GetUserInfoResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "public: GetUserInfo")
	defer span.End()

	userData, err := auth.GetUserData(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserData: %w", err)
	}

	user, err := c.svc.GetUserInfo(ctx, userData.UserID)
	if err != nil {
		return nil, fmt.Errorf("svc.GetUserInfo: %w", err)
	}

	return &public.GetUserInfoResponse{
		Id:    uint32(userData.UserID),
		Email: user.Email,
		Role:  enum.ConvertServiceRoleToPublic(user.Role),
	}, nil
}
