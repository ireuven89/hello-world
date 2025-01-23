package biddering

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
	"github.com/labstack/gommon/log"

	"github.com/ireuven89/hello-world/backend/biddering/model"
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
	log.Printf("Starting server on port %s...", port)
	err := http.ListenAndServe(":"+port, t.router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func RegisterRoutes(router *httprouter.Router, s Service) {

	getBidderHandler := kithttp.NewServer(
		MakeEndpointGetBidder(s),
		decodeGetBidderRequest,
		encodeGetBidderResponse,
	)

	listBiddersHandler := kithttp.NewServer(
		MakeEndpointListBidders(s),
		decodeListBiddersRequest,
		encodeListBiddersResponse,
	)

	createBidderHandler := kithttp.NewServer(
		MakeEndpointCreateBidder(s),
		decodeCreateBidderRequest,
		encodeCreateBidderResponse,
	)

	updateBidderHandler := kithttp.NewServer(
		MakeEndpointUpdateBidder(s),
		decodeUpdateBidderRequest,
		kithttp.EncodeJSONResponse,
	)

	deleteBidderHandler := kithttp.NewServer(
		MakeEndpointDeleteBidder(s),
		decodeDeleteBidderRequest,
		kithttp.EncodeJSONResponse,
	)

	router.Handler(http.MethodGet, "/bidders/:uuid", getBidderHandler)
	router.Handler(http.MethodGet, "/bidders", listBiddersHandler)
	router.Handler(http.MethodPost, "/bidders", createBidderHandler)
	router.Handler(http.MethodPut, "/bidders", updateBidderHandler)
	router.Handler(http.MethodDelete, "/bidders", deleteBidderHandler)

}

func decodeGetBidderRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return GetBidderRequest{
		uuid: r.PathValue("uuid"),
	}, nil
}

func encodeGetBidderResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	result := response.(GetBidderResponse)

	//todo format result

	return json.NewEncoder(w).Encode(result)
}

func decodeListBiddersRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	params := r.URL.Query()

	return ListBiddersRequest{
		Input: model.BiddersInput{
			Uuid: params.Get("uuid"),
			Name: params.Get("name"),
			Item: params.Get("item"),
			Page: model.SetPage("offset", "limit"),
		},
	}, nil
}

func encodeListBiddersResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	result := response.(GetBidderResponse)

	//todo format result

	return json.NewEncoder(w).Encode(result)
}

func decodeCreateBidderRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req CreateBidderRequestModel

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("transport.decodeCreateBidderRequest failed decoding request")
		return nil, err
	}

	return req, nil
}

func encodeCreateBidderResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	result, ok := response.(CreateBidderResponseModel)

	if !ok {
		return errors.New("transport.encodeCreateBidderResponse failed encoding response")
	}

	//todo format result

	return json.NewEncoder(w).Encode(result)
}

func decodeUpdateBidderRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req UpdateBidderRequestModel

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("transport.decodeUpdateBidderRequest failed decoding request")
		return nil, err
	}

	return req, nil
}

func decodeDeleteBidderRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req DeleteBidderRequestModel

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("transport.decodeUpdateBidderRequest failed decoding request")
		return nil, err
	}

	return req, nil
}
