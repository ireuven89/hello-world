package biddering

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"

	"github.com/ireuven89/hello-world/backend/biddering/model"
)

type GetBidderRequest struct {
	uuid string
}

type GetBidderResponse struct {
	model.Bidder
}

func MakeEndpointGetBidder(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetBidderRequest)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointGetUser failed cast request")
		}

		result, err := s.GetBidder(req.uuid)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointRegister: %v", err)
		}

		return result, nil
	}
}

type ListBiddersRequest struct {
	Input model.BiddersInput
}

type ListBiddersResponseModel struct {
	bidders []model.Bidder
}

func MakeEndpointListBidders(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(ListBiddersRequest)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointListBidders failed cast request")
		}

		result, err := s.ListBidders(req.Input)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointListBidders: %v", err)
		}

		return ListBiddersResponseModel{
			bidders: result,
		}, nil
	}
}

type CreateBidderRequestModel struct {
	Input model.BidderInput
}

type CreateBidderResponseModel struct {
	ID string `json:"id"`
}

func MakeEndpointCreateBidder(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(CreateBidderRequestModel)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointCreateBidder failed cast request")
		}

		id, err := s.CreateBidder(req.Input)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointCreateBidder: %v", err)
		}

		return CreateBidderResponseModel{
			ID: id,
		}, nil
	}
}

type UpdateBidderRequestModel struct {
	Input model.BidderInput
}

func MakeEndpointUpdateBidder(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(UpdateBidderRequestModel)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointCreateBidder failed cast request")
		}

		_, err = s.UpdateBidder(req.Input)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointCreateBidder: %v", err)
		}

		return nil, nil
	}
}

type DeleteBidderRequestModel struct {
	Uuid string
}

func MakeEndpointDeleteBidder(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(DeleteBidderRequestModel)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointCreateBidder failed cast request")
		}

		err = s.Delete(req.Uuid)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointCreateBidder: %v", err)
		}

		return nil, nil
	}
}
