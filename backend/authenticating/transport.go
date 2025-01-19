package authenticating

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
	"github.com/labstack/gommon/log"
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
<<<<<<< HEAD
	log.Printf("Starting auth server on port %s...", port)
=======
	log.Printf("Starting server on port %s...", port)
>>>>>>> 4f75e37 (authenticathing service and additional changes)
	err := http.ListenAndServe(":"+port, t.router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func RegisterRoutes(router *httprouter.Router, s Service) {

<<<<<<< HEAD
	healthHandler := kithttp.NewServer(
		MakeEndpointHealth(s),
		decodeHealthRequest,
		kithttp.EncodeJSONResponse,
	)

=======
>>>>>>> 4f75e37 (authenticathing service and additional changes)
	registerUserHandler := kithttp.NewServer(
		MakeEndpointRegister(s),
		decodeRegisterRequest,
		kithttp.EncodeJSONResponse,
	)

	loginUserHandler := kithttp.NewServer(
		MakeEndpointLogin(s),
		decodeLoginRequest,
		encodeLoginResponse,
	)

	verifyTokenHandler := kithttp.NewServer(
		MakeEndpointVerify(s),
		decodeVerifyRequest,
		encodeVerifyResponse,
	)

<<<<<<< HEAD
	router.Handler(http.MethodGet, "/health", healthHandler)
=======
>>>>>>> 4f75e37 (authenticathing service and additional changes)
	router.Handler(http.MethodPost, "/register", registerUserHandler)
	router.Handler(http.MethodPost, "/login", loginUserHandler)
	router.Handler(http.MethodPost, "/verify", verifyTokenHandler)
}

<<<<<<< HEAD
func decodeHealthRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return nil, nil
}

=======
>>>>>>> 4f75e37 (authenticathing service and additional changes)
func decodeRegisterRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req RegisterRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeLoginRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req LoginRequestModel

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func encodeLoginResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	res, ok := response.(LoginResponseModel)

	if !ok {
		return errors.New("encodeLoginResponse failed to decode response")
	}

	formatted := map[string]interface{}{
		"token": res.Token,
	}

	writer.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(writer).Encode(formatted)
}

func decodeVerifyRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req VerifyRequestModel

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func encodeVerifyResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	res, ok := response.(VerifyResponseModel)

	if !ok {
		return errors.New("encodeVerifyResponse failed to decode response")
	}

	formatted := map[string]interface{}{
		"user": res.User,
	}

	writer.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(writer).Encode(formatted)
}
