package users

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
	"github.com/labstack/gommon/log"

	"github.com/ireuven89/hello-world/backend/users/model"
)

type Router interface {
	Handle(method, path string, handler http.Handler)
}

func NewTransport(s Service, router *httprouter.Router) Transport {

	transport := Transport{
		router: router,
		s:      s,
	}
	RegisterRoutes(router, s) // Register routes during initialization
	return transport
}

type Transport struct {
	router *httprouter.Router
	s      Service
}

func (t *Transport) ListenAndServe(port string) {
	log.Printf("Starting server on port %s...", port)
	err := http.ListenAndServe(":"+port, t.router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func RegisterRoutes(router *httprouter.Router, s Service) {
	getUserHandler := kithttp.NewServer(
		MakeEndpointGetUser(s),
		decodeGetUserRequest,
		encodeGetUserResponse,
	)

	getUsersHandler := kithttp.NewServer(
		MakeEndpointGetUsers(s),
		decodeListUserRequest,
		encodeListUsersResponse,
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
	router.Handler(http.MethodGet, "/users", getUsersHandler)
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

type ListUserRequest struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

type ListUserResponse struct {
	users []model.UserResponse
}

func decodeListUserRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req ListUserRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func encodeListUsersResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	res, ok := response.(ListUserResponse)
	if !ok {
		return fmt.Errorf("encodeListUsersResponse failed cast response")
	}

	return json.NewEncoder(writer).Encode(&res)
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
