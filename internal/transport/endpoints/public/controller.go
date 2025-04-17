package public

import (
	"github.com/go-kit/kit/endpoint"
	"go.uber.org/zap"

	"users-service/internal/users"

	public "github.com/knstch/users-api/public"
)

type Endpoints struct {
	CreateUser endpoint.Endpoint
}

type Controller struct {
	svc users.Users
	lg  *zap.Logger

	public.UnimplementedUsersServer
}

func MakeEndpoints(c *Controller) Endpoints {
	return Endpoints{
		CreateUser: MakeCreateUserEndpoint(c),
	}
}

func NewController(svc users.Users, lg *zap.Logger) *Controller {
	return &Controller{
		svc: svc,
		lg:  lg,
	}
}
