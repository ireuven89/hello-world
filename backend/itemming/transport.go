package itemming

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
	"github.com/labstack/gommon/log"

	"github.com/ireuven89/hello-world/backend/itemming/model"
)

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

type Router interface {
	Handle(method, path string, handler http.Handler)
}

func (t *Transport) ListenAndServe(port string) {
	log.Printf("Starting item server on port %s...", port)
	err := http.ListenAndServe(":"+port, t.router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func RegisterRoutes(router *httprouter.Router, s Service) {

	healthHandler := kithttp.NewServer(
		MakeEndpointHealth(s),
		decodeHealthItemRequest,
		kithttp.EncodeJSONResponse,
	)

	getItemHandler := kithttp.NewServer(
		MakeEndpointGetItem(s),
		decodeGetItemRequest,
		encodeGetItemResponse,
	)

	listItemsHandler := kithttp.NewServer(
		MakeEndpointListItems(s),
		decodeListItemsRequest,
		encodeListItemsResponse,
	)

	createItemHandler := kithttp.NewServer(
		MakeEndpointCreateItem(s),
		decodeCreateItemRequest,
		encodeCreateItemResponse,
	)

	createItemsHandler := kithttp.NewServer(
		MakeEndpointCreateItems(s),
		decodeCreateItemRequest,
		kithttp.EncodeJSONResponse,
	)

	updateItemHandler := kithttp.NewServer(
		MakeEndpointUpdateItem(s),
		decodeUpdateItemRequest,
		kithttp.EncodeJSONResponse,
	)

	deleteItemHandler := kithttp.NewServer(
		MakeEndpointDeleteItem(s),
		decodeDeleteItemRequest,
		kithttp.EncodeJSONResponse,
	)

	router.Handler(http.MethodGet, "/health", healthHandler)
	router.Handler(http.MethodGet, "/items/:uuid", getItemHandler)
	router.Handler(http.MethodGet, "/items", listItemsHandler)
	router.Handler(http.MethodPost, "/items", createItemHandler)
	router.Handler(http.MethodPut, "/items", createItemsHandler)
	router.Handler(http.MethodPatch, "/items/:uuid", updateItemHandler)
	router.Handler(http.MethodDelete, "/items/:uuid", deleteItemHandler)
}

func decodeHealthItemRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {

	return nil, nil
}

func decodeGetItemRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {

	return GetItemRequest{
		Uuid: r.PathValue("uuid"),
	}, nil
}

func encodeGetItemResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	res, ok := response.(GetItemResponse)

	if !ok {
		return errors.New("encodeGetItemResponse.failed encode response")
	}

	return json.NewEncoder(writer).Encode(res)
}

func decodeListItemsRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	queryParams := r.URL.Query()
	name := queryParams.Get("name")
	description := queryParams.Get("description")
	price := queryParams.Get("link")

	return ListItemsRequest{
		input: model.ListInput{
			Price:       price,
			Description: description,
			Name:        name,
		},
	}, nil
}

func encodeListItemsResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	res, ok := response.(ListItemsResponse)

	if !ok {
		return errors.New("encodeListItemsResponse failed to parse response")
	}

	formatted := map[string]interface{}{
		"items": res.items,
	}

	writer.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(writer).Encode(formatted)
}

func decodeCreateItemRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req CreateItemRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func encodeCreateItemResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	res, ok := response.(CreateItemResponse)

	if !ok {
		return errors.New("encodeCreateItemResponse failed to encode response")
	}

	formatted := map[string]interface{}{
		"id": res.Uuid,
	}

	writer.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(writer).Encode(formatted)
}

func decodeUpdateItemRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req UpdateItemRequest

	req.item.Uuid = r.PathValue("uuid")

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeDeleteItemRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {

	return GetItemRequest{
		Uuid: r.PathValue("uuid"),
	}, nil
}
