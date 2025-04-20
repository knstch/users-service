package public

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	public "github.com/knstch/users-api/public"
)

func MakeCreateUserEndpoint(c *Controller) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return c.CreateUser(ctx, request.(*public.CreateUserRequest))
	}
}

func (c *Controller) CreateUser(ctx context.Context, req *public.CreateUserRequest) (*public.CreateUserResponse, error) {
	tokens, err := c.svc.CreateUser(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &public.CreateUserResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
