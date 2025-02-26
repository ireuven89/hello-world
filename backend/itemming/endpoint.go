package itemming

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"

	"github.com/ireuven89/hello-world/backend/itemming/model"
)

type GetItemRequest struct {
	Uuid string
}

type GetItemResponse struct {
	item model.Item
}

func MakeEndpointGetItem(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetItemRequest)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointGetItem failed cast request")
		}

		result, err := s.GetItem(req.Uuid)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointGetItem: %v", err)
		}

		return result, nil
	}
}

type ListItemsRequest struct {
	input model.ListInput
}

type ListItemsResponse struct {
	items []model.Item
}

func MakeEndpointListItems(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(ListItemsRequest)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointGetUser failed cast request")
		}

		result, err := s.GetItems(req.input)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointGetUser: %v", err)
		}

		return ListItemsResponse{
			items: result,
		}, nil
	}
}

type CreateItemRequest struct {
	item model.ItemInput
}

type CreateItemResponse struct {
	Uuid string `json:"uuid"`
}

func MakeEndpointCreateItem(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(CreateItemRequest)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointGetUser failed cast request")
		}

		id, err := s.CreateItem(req.item)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointGetUser: %v", err)
		}

		return CreateItemResponse{
			Uuid: id,
		}, nil
	}
}

type UpdateItemRequest struct {
	item model.ItemInput
}

func MakeEndpointUpdateItem(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(UpdateItemRequest)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointGetUser failed cast request")
		}

		err = s.UpdateItem(req.item)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointUpdateItem: %v", err)
		}

		return nil, nil
	}
}

type DeleteItemRequest struct {
	Uuid string
}

type DeleteItemResponse struct {
	item model.Item
}

func MakeEndpointDeleteItem(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetItemRequest)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointGetItem failed cast request")
		}

		err = s.DeleteItem(req.Uuid)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointGetItem: %v", err)
		}

		return nil, nil
	}
}
