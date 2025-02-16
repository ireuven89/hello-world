package authenticating

import (
	"context"
	"fmt"

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

func MakeEndpointHealth(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		return s.Health(), nil
	}
}

func MakeEndpointRegister(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(RegisterRequest)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointRegister failed cast request")
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
			return nil, fmt.Errorf("MakeEndpointLogin failed cast request")
		}

		token, err := s.Login(req.UserName, req.Password)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointLogin: %v", err)
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
			return nil, fmt.Errorf("MakeEndpointVerify failed cast request")
		}

		user, err := s.VerifyToken(req.JwtToken)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointVerify: %v", err)
		}

		return VerifyResponseModel{
			User: user,
		}, nil
	}
}
