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
			Path:    "/createUser",
			Handler: MakeCreateUserEndpoint(c),
			Decoder: transport.DecodeJSONRequest[public.CreateUserRequest],
			Encoder: httptransport.EncodeJSONResponse,
			Req:     public.CreateUserRequest{},
			Res:     public.CreateUserResponse{},
		},
		{
			Method:  http.MethodPost,
			Path:    "/confirmEmail",
			Handler: MakeConfirmEmailEndpoint(c),
			Decoder: transport.DecodeJSONRequest[public.ConfirmEmailRequest],
			Encoder: httptransport.EncodeJSONResponse,
			Req:     public.ConfirmEmailRequest{},
			Res:     public.ConfirmEmailResponse{},
			Mdw:     mdw,
		},
		{
			Method:  http.MethodPost,
			Path:    "/refreshToken",
			Handler: MakeRefreshTokenEndpoint(c),
			Decoder: transport.DecodeJSONRequest[public.RefreshTokenRequest],
			Encoder: httptransport.EncodeJSONResponse,
			Req:     public.RefreshTokenRequest{},
			Res:     public.RefreshTokenResponse{},
		},
	}
}
