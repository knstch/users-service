package transport

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"users-service/internal/transport/endpoints/public"

	publicApi "github.com/knstch/users-api/public"
)

func DecodeJSONRequest[T any](_ context.Context, r *http.Request) (interface{}, error) {
	var req T
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func EncodeJSONResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func NewHTTPHandlers(endpoints public.Endpoints) http.Handler {
	r := mux.NewRouter()

	r.Methods("POST").Path("/createUser").Handler(httptransport.NewServer(
		endpoints.CreateUser,
		DecodeJSONRequest[publicApi.CreateUserRequest],
		EncodeJSONResponse,
	))

	return r
}
