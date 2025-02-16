package auctioning

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"

	"github.com/ireuven89/hello-world/backend/auctioning/model"
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
	log.Printf("Starting auction server on port %s...", port)
	err := http.ListenAndServe(":"+port, t.router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func RegisterRoutes(router *httprouter.Router, s Service) {

	healthHandler := kithttp.NewServer(
		MakeEndpointHealth(s),
		decodeHealthRequest,
		kithttp.EncodeJSONResponse,
	)

	getAuctionHandler := kithttp.NewServer(
		MakeEndpointGet(s),
		decodeGetRequest,
		encodeGetResponse,
	)

	listAuctionHandler := kithttp.NewServer(
		MakeEndpointSearch(s),
		decodeFetchRequest,
		encodeFetchResponse,
	)

	deleteAuctionHandler := kithttp.NewServer(
		MakeEndpointDeleteAuction(s),
		decodeDeleteRequest,
		kithttp.EncodeJSONResponse,
	)

	updateAuctionHandler := kithttp.NewServer(
		MakeEndpointUpdateAuction(s),
		decodeUpdateRequest,
		kithttp.EncodeJSONResponse,
	)

	router.Handler(http.MethodGet, "/health", healthHandler)
	router.Handler(http.MethodGet, "/auctioning/:uuid", getAuctionHandler)
	router.Handler(http.MethodPost, "/auctioning", listAuctionHandler)
	router.Handler(http.MethodDelete, "/auctioning/:uuid", deleteAuctionHandler)
	router.Handler(http.MethodPut, "/auctioning/:uuid", updateAuctionHandler)
}

func decodeHealthRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return nil, nil
}

func decodeGetRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return GetRequest{
		Uuid: r.PathValue("uuid"),
	}, nil
}

func encodeGetResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	res, ok := response.(GetResponse)

	if !ok {
		return fmt.Errorf("Transport.encodeGetResponse failed encoding response ")
	}

	return json.NewEncoder(writer).Encode(res)
}

func decodeFetchRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req FetchRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func encodeFetchResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	res, ok := response.(FetchAllResponse)

	if !ok {
		return fmt.Errorf("encodeFetchResponse failed to decode response")
	}

	formatted := map[string]interface{}{
		"result": res.auctions,
	}

	writer.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(writer).Encode(formatted)
}

func decodeDeleteRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {

	return DeleteRequest{
		Uuid: r.PathValue("uuid"),
	}, nil
}

func decodeUpdateRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.AuctionRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	req.Id = r.PathValue("uuid")

	return req, nil
}
