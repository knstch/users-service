package public

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/knstch/subtrack-libs/tracing"
	public "github.com/knstch/users-api/public"
)

func MakeRegisterEndpoint(c *Controller) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return c.Register(ctx, request.(*public.RegisterRequest))
	}
}

func (c *Controller) Register(ctx context.Context, req *public.RegisterRequest) (*public.RegisterResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "public: Register")
	defer span.End()

	tokens, err := c.svc.Register(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &public.RegisterResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
