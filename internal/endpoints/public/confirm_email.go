package public

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/knstch/subtrack-libs/tracing"
	public "github.com/knstch/users-api/public"
)

func MakeConfirmEmailEndpoint(c *Controller) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return c.ConfirmEmail(ctx, request.(*public.ConfirmEmailRequest))
	}
}

func (c *Controller) ConfirmEmail(ctx context.Context, req *public.ConfirmEmailRequest) (*public.ConfirmEmailResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "public: ConfirmEmail")
	defer span.End()

	tokens, err := c.svc.ConfirmEmail(ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("svc.ConfirmEmail: %w", err)
	}

	return &public.ConfirmEmailResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
