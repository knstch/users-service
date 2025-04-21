package public

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	public "github.com/knstch/users-api/public"
)

func MakeRefreshTokenEndpoint(c *Controller) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return c.RefreshToken(ctx, request.(*public.RefreshTokenRequest))
	}
}

func (c *Controller) RefreshToken(ctx context.Context, req *public.RefreshTokenRequest) (*public.RefreshTokenResponse, error) {
	tokens, err := c.svc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("svc.RefreshToken: %w", err)
	}

	return &public.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
