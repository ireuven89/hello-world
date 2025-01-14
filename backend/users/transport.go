package users

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/ireuven89/hello-world/backend/users/model"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Router interface {
	Handle(method, path string, handler http.Handler)
}

type Transport struct {
	router *httprouter.Router
	s      Service
}

func NewTransport(router *httprouter.Router, s Service) Transport {
	return Transport{
		router: router,
		s:      s,
	}
}

func RegisterRoutes(router *httprouter.Router, s Service) {
	getUserHandler := kithttp.NewServer(
		MakeEndpointGetUser(s),
		decodeGetUserRequest,
		encodeGetUserResponse,
	)

	createUserHandler := kithttp.NewServer(
		MakeEndpointCreateUser(s),
		decodeCreateUserRequest,
		encodeCreateUserResponse,
	)

	updateUserHandler := kithttp.NewServer(
		MakeEndpointUpdateUser(s),
		decodeDeleteUserRequest,
		encodeDeleteUserResponse,
	)

	deleteUserHandler := kithttp.NewServer(
		MakeEndpointDeleteUser(s),
		decodeDeleteUserRequest,
		encodeDeleteUserResponse,
	)

	router.Handler(http.MethodGet, "/users/:id", getUserHandler)
	router.Handler(http.MethodPost, "/users", createUserHandler)
	router.Handler(http.MethodPut, "/users", updateUserHandler)
	router.Handler(http.MethodDelete, "/users/:id", deleteUserHandler)
}

type GetUserRequest struct {
	uuid string
}

type GetUserResponse struct {
	model.UserResponse
}

func decodeGetUserRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	params := httprouter.ParamsFromContext(ctx)
	uuid := params.ByName("uuid")

	if uuid == "" {
		return nil, errors.New("invalid param")
	}

	return GetUserRequest{
		uuid: uuid,
	}, nil
}

func encodeGetUserResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	_, ok := response.(GetUserResponse)
	if !ok {
		return fmt.Errorf("encodeGetArticleResponse failed cast response")
	}

	return nil
}

type CreateUserRequest struct {
	model.UserUpsertInput
}

type CreateUserResponse struct {
	model.UserUpsertInput
}

func decodeCreateUserRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req CreateUserRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func encodeCreateUserResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	res, ok := response.(model.UserResponse)

	if !ok {
		return fmt.Errorf("encodeGetUserResponse failed cast response")
	}

	if err := json.NewEncoder(writer).Encode(res); err != nil {
		return fmt.Errorf("encodeGetUserResponse failed cast response")
	}

	return nil
}

type DeleteUserRequest struct {
	uuid string
}

func decodeDeleteUserRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req DeleteUserRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func encodeDeleteUserResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {

	return nil
}
