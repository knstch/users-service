package public

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/knstch/subtrack-libs/middleware"
	"go.uber.org/zap"

	"users-service/config"
	"users-service/internal/users"

	public "github.com/knstch/users-api/public"

	"github.com/knstch/subtrack-libs/endpoints"
	"github.com/knstch/subtrack-libs/transport"
)

type Endpoints struct {
	CreateUser endpoint.Endpoint
}

type Controller struct {
	svc users.Users
	lg  *zap.Logger
	cfg *config.Config

	public.UnimplementedUsersServer
}

func NewController(svc users.Users, lg *zap.Logger, cfg *config.Config) *Controller {
	return &Controller{
		svc: svc,
		cfg: cfg,
		lg:  lg,
	}
}

func (c *Controller) Endpoints() []endpoints.Endpoint {
	mdw := []middleware.Middleware{middleware.WithCookieAuth(c.cfg.JwtSecret)}

	return []endpoints.Endpoint{
		{
			Method:  http.MethodPost,
			Path:    "/register",
			Handler: MakeRegisterEndpoint(c),
			Decoder: transport.DecodeJSONRequest[public.RegisterRequest],
			Encoder: httptransport.EncodeJSONResponse,
		},
		{
			Method:  http.MethodPost,
			Path:    "/confirmEmail",
			Handler: MakeConfirmEmailEndpoint(c),
			Decoder: transport.DecodeJSONRequest[public.ConfirmEmailRequest],
			Encoder: httptransport.EncodeJSONResponse,
			Mdw:     mdw,
		},
		{
			Method:  http.MethodPost,
			Path:    "/refreshToken",
			Handler: MakeRefreshTokenEndpoint(c),
			Decoder: transport.DecodeJSONRequest[public.RefreshTokenRequest],
			Encoder: httptransport.EncodeJSONResponse,
		},
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: MakeLoginEndpoint(c),
			Decoder: transport.DecodeJSONRequest[public.LoginRequest],
			Encoder: httptransport.EncodeJSONResponse,
		},
		{
			Method:  http.MethodGet,
			Path:    "/getUserInfo",
			Handler: MakeGetUserInfoEndpoint(c),
			Decoder: transport.DecodeDefaultRequest,
			Encoder: httptransport.EncodeJSONResponse,
			Mdw:     mdw,
		},
		{
			Method:  http.MethodPost,
			Path:    "/resetPassword",
			Handler: MakeResetPasswordEndpoint(c),
			Decoder: transport.DecodeJSONRequest[public.ResetPasswordRequest],
			Encoder: httptransport.EncodeJSONResponse,
		},
		{
			Method:  http.MethodPost,
			Path:    "/confirmResetPassword",
			Handler: MakeConfirmResetPasswordEndpoint(c),
			Decoder: transport.DecodeJSONRequest[public.ConfirmResetPasswordRequest],
			Encoder: httptransport.EncodeJSONResponse,
		},
	}
}
