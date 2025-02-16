package userring

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"

	"github.com/ireuven89/hello-world/backend/userring/model"
)

func MakeEndpointGetUser(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(string)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointGetUser failed cast request")
		}

		user, err := s.GetUser(req)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointGetUser: %v", err)
		}

		return model.UserResponse{
			Name:        user.Name,
			Uuid:        user.Name,
			Description: user.Name,
		}, nil
	}
}

func MakeEndpointGetUsers(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var result []model.UserResponse
		req, ok := request.(model.UserFetchInput)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointGetUser failed cast request")
		}

		users, err := s.ListUsers(req)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointGetUser: %v", err)
		}

		for _, user := range users {
			result = append(result, model.UserResponse{
				Name:        user.Name,
				Uuid:        user.Name,
				Description: user.Name,
			})
		}

		return result, nil
	}
}

func MakeEndpointCreateUser(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(string)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointGetUser failed cast request")
		}

		user, err := s.GetUser(req)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointGetUser: %v", err)
		}

		return model.UserResponse{
			Name:        user.Name,
			Uuid:        user.Name,
			Description: user.Name,
		}, nil
	}
}

func MakeEndpointUpdateUser(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(string)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointGetUser failed cast request")
		}

		user, err := s.GetUser(req)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointGetUser: %v", err)
		}

		return model.UserResponse{
			Name:        user.Name,
			Uuid:        user.Name,
			Description: user.Name,
		}, nil
	}
}

func MakeEndpointDeleteUser(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(string)
		if !ok {
			return "", fmt.Errorf("MakeEndpointGetUser failed cast request")
		}

		if err = s.DeleteUser(req); err != nil {
			return nil, err
		}

		return "", nil
	}
}

func MakeEndpointMakeUser(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(string)

		if !ok {
			return "", fmt.Errorf("MakeEndpointGetUser failed to ")
		}

		return req, nil
	}
}
