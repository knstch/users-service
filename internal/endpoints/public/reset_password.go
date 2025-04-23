package public

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	public "github.com/knstch/users-api/public"
)

func MakeResetPasswordEndpoint(c *Controller) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return c.ResetPassword(ctx, request.(*public.ResetPasswordRequest))
	}
}

func (c *Controller) ResetPassword(ctx context.Context, req *public.ResetPasswordRequest) (*public.ResetPasswordResponse, error) {
	if err := c.svc.ResetPassword(ctx, req.Email); err != nil {
		return nil, fmt.Errorf("svc.ResetPassword: %w", err)
	}

	return &public.ResetPasswordResponse{}, nil
}
