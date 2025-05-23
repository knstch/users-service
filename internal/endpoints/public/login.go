package public

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/knstch/subtrack-libs/tracing"
	public "github.com/knstch/users-api/public"
)

func MakeLoginEndpoint(c *Controller) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return c.Login(ctx, request.(*public.LoginRequest))
	}
}

func (c *Controller) Login(ctx context.Context, req *public.LoginRequest) (*public.LoginResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "public: Login")
	defer span.End()

	tokens, err := c.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("svc.Login: %w", err)
	}

	return &public.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
