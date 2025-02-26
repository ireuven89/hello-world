package authenticating

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
)

type RegisterRequest struct {
	UserName string `json:"user"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	UserName string `json:"model"`
	Password string `json:"password"`
}

func MakeEndpointRegister(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(RegisterRequest)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointGetUser failed cast request")
		}

		err = s.Register(req.UserName, req.Password)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointRegister: %v", err)
		}

		return nil, nil
	}
}

type LoginRequestModel struct {
	UserName string `json:"user"`
	Password string `json:"password"`
}

type LoginResponseModel struct {
	Token string `json:"token"`
}

func MakeEndpointLogin(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(LoginRequestModel)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointGetUser failed cast request")
		}

		token, err := s.Login(req.UserName, req.Password)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointGetUser: %v", err)
		}

		return LoginResponseModel{
			Token: token,
		}, nil
	}
}

type VerifyRequestModel struct {
	JwtToken string `json:"JwtToken"`
}

type VerifyResponseModel struct {
	User string `json:"user"`
}

func MakeEndpointVerify(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(VerifyRequestModel)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointGetUser failed cast request")
		}

		user, err := s.VerifyToken(req.JwtToken)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointGetUser: %v", err)
		}

		return VerifyResponseModel{
			User: user,
		}, nil
	}
}

type MigrateRequest struct {
	time time.Time
}

func MakeEndpointMigrate(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(time.Time)
		err = s.Migrate(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointGetUser: %v", err)
		}

		return nil, nil
	}
}
