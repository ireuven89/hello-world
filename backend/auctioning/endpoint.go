package auctioning

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"

	"github.com/ireuven89/hello-world/backend/auctioning/model"
)

func MakeEndpointHealth(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		return s.Health(), nil
	}
}

type FetchRequest struct {
	req model.AuctionRequest
}

type FetchAllResponse struct {
	auctions []model.Auction
}

func MakeEndpointSearch(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(FetchRequest)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointRegister failed cast request")
		}

		result, err := s.Search(req.req)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointRegister: %v", err)
		}

		return result, nil
	}
}

type GetRequest struct {
	Uuid string `json:"uuid"`
}

type GetResponse struct {
	Auction model.Auction
}

func MakeEndpointGet(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetRequest)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointLogin failed cast request")
		}

		auction, err := s.Find(req.Uuid)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointLogin: %v", err)
		}

		return GetResponse{
			Auction: auction,
		}, nil
	}
}

type DeleteRequest struct {
	Uuid string `json:"uuid"`
}

func MakeEndpointDeleteAuction(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(DeleteRequest)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointLogin failed cast request")
		}

		err = s.Delete(req.Uuid)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointLogin: %v", err)
		}

		return nil, nil
	}
}

type UpdateRequest struct {
	req model.AuctionRequest
}

func MakeEndpointUpdateAuction(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(UpdateRequest)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointLogin failed cast request")
		}

		err = s.Update(req.req)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointLogin: %v", err)
		}

		return nil, nil
	}
}
