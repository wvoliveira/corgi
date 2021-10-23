package auth

// The Auth is just over HTTP, so we just have a single transport.go.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"

	"github.com/heptiolabs/healthcheck"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
// Useful in a auth service server.
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	postSignupHandler := kithttp.NewServer(
		e.PostSignupEndpoint,
		decodePostLoginRequest,
		encodeResponse,
		options...,
	)

	postLoginHandler := kithttp.NewServer(
		e.PostLoginEndpoint,
		decodePostLoginRequest,
		encodeResponse,
		options...,
	)

	postLogoutHandler := kithttp.NewServer(
		e.PostLogoutEndpoint,
		decodePostLogoutRequest,
		encodeResponse,
		options...,
	)

	postRefreshHandler := kithttp.NewServer(
		e.PostRefreshEndpoint,
		decodePostRefreshRequest,
		encodeResponse,
		options...,
	)

	health := healthcheck.NewHandler()
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))
	health.AddReadinessCheck("upstream-dep-dns", healthcheck.DNSResolveCheck("localhost", 50*time.Millisecond))

	r.Handle("/auth/v1/signup", postSignupHandler).Methods("POST")
	r.Handle("/auth/v1/login", postLoginHandler).Methods("POST")
	r.Handle("/auth/v1/logout", postLogoutHandler).Methods("POST")
	r.Handle("/auth/v1/refresh", postRefreshHandler).Methods("POST")

	r.HandleFunc("/auth/v1/health/ready", health.ReadyEndpoint)
	r.HandleFunc("/auth/v1/health/live", health.LiveEndpoint)

	r.PathPrefix("/auth/v1/swagger").Handler(httpSwagger.WrapHandler)

	return r
}

// PostSignup godoc
// @Summary Create new JWT
// @Description Create new JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body auth.PostSignup true "Signup struct"
// @Success 200
// @Router /auth/v1/signup [post]
func decodePostSignupRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postSignupRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Signup); e != nil {
		return nil, e
	}
	return req, nil
}

// PostLogin godoc
// @Summary Create new JWT
// @Description Create new JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body auth.PostLogin true "Login struct"
// @Success 200
// @Router /auth/v1/login [post]
func decodePostLoginRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postLoginRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Login); e != nil {
		return nil, e
	}
	return req, nil
}

// PostLogout godoc
// @Summary Delete JWT
// @Description Delete JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body auth.PostLogout true "Logout struct"
// @Success 200
// @Router /auth/v1/logout [post]
func decodePostLogoutRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postLogoutRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Logout); e != nil {
		return nil, e
	}
	return req, nil
}

// PostRefresh godoc
// @Summary Refresh JWT
// @Description Refresh JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body auth.PostRefresh true "Refresh struct"
// @Success 200
// @Router /auth/v1/refresh [post]
func decodePostRefreshRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postRefreshRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Refresh); e != nil {
		return nil, e
	}
	return req, nil
}

func encodePostSignupRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("POST").Path("/signup")
	req.URL.Path = "/auth/v1/signup"
	return encodeRequest(ctx, req, request)
}

func encodePostLoginRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("POST").Path("/login")
	req.URL.Path = "/auth/v1/login"
	return encodeRequest(ctx, req, request)
}

func encodePostLogoutRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("POST").Path("/logout")
	req.URL.Path = "/auth/v1/logout"
	return encodeRequest(ctx, req, request)
}

func encodePostRfreshRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("POST").Path("/refresh")
	req.URL.Path = "/auth/v1/refresh"
	return encodeRequest(ctx, req, request)
}

func decodePostSignupResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response postSignupResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePostLoginResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response postLoginResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePostLogoutResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response postLogoutResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}
func decodePostRefreshResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response postRefreshResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// auth service endpoints require mutating the HTTP method and request path.
func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists, ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
