package public

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/knstch/subtrack-libs/tracing"
	public "github.com/knstch/users-api/public"
)

func MakeConfirmResetPasswordEndpoint(c *Controller) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return c.ConfirmResetPassword(ctx, request.(*public.ConfirmResetPasswordRequest))
	}
}

func (c *Controller) ConfirmResetPassword(ctx context.Context, req *public.ConfirmResetPasswordRequest) (*public.ConfirmResetPasswordResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "public: ConfirmResetPassword")
	defer span.End()

	if err := c.svc.ConfirmResetPassword(ctx, req.Email, req.Code, req.Password); err != nil {
		return nil, fmt.Errorf("svc.ConfirmResetPassword: %w", err)
	}

	return &public.ConfirmResetPasswordResponse{}, nil
}
