package user

// The "user" is just over HTTP, so we just have a single transport.go.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
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
// Useful in a "user" service.
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	getUsersHandler := kithttp.NewServer(
		e.GetUsersEndpoint,
		decodeGetUsersRequest,
		encodeResponse,
		options...,
	)

	postUserHandler := kithttp.NewServer(
		e.PostUserEndpoint,
		decodePostUserRequest,
		encodeResponse,
		options...,
	)

	getUserHandler := kithttp.NewServer(
		e.GetUserEndpoint,
		decodeGetUserRequest,
		encodeResponse,
		options...,
	)

	putUserHandler := kithttp.NewServer(
		e.PutUserEndpoint,
		decodePutUserRequest,
		encodeResponse,
		options...,
	)

	patchUserHandler := kithttp.NewServer(
		e.PatchUserEndpoint,
		decodePatchUserRequest,
		encodeResponse,
		options...,
	)

	deleteUserHandler := kithttp.NewServer(
		e.DeleteUserEndpoint,
		decodeDeleteUserRequest,
		encodeResponse,
		options...,
	)

	health := healthcheck.NewHandler()
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))
	health.AddReadinessCheck("upstream-dep-dns", healthcheck.DNSResolveCheck("localhost", 50*time.Millisecond))

	r.Handle("/user/v1/users", getUsersHandler).Methods("GET")
	r.Handle("/user/v1/users", postUserHandler).Methods("POST")
	r.Handle("/user/v1/users/{id}", getUserHandler).Methods("GET")
	r.Handle("/user/v1/users/{id}", putUserHandler).Methods("PUT")
	r.Handle("/user/v1/users/{id}", patchUserHandler).Methods("PATCH")
	r.Handle("/user/v1/users/{id}", deleteUserHandler).Methods("DELETE")

	r.HandleFunc("/user/v1/health/ready", health.ReadyEndpoint)
	r.HandleFunc("/user/v1/health/live", health.LiveEndpoint)

	r.PathPrefix("/user/v1/swagger").Handler(httpSwagger.WrapHandler)

	return r
}

// PostUser godoc
// @Summary Add a new User
// @Description Add a new User
// @Tags User
// @Accept json
// @Produce json
// @Param data body user.PostUser true "User struct"
// @Success 200
// @Router /user/v1/users [post]
func decodePostUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postUserRequest
	if e := json.NewDecoder(r.Body).Decode(&req.User); e != nil {
		return nil, e
	}
	return req, nil
}

// GetUser godoc
// @Summary Get an User
// @Description Get User by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200
// @Failure 404
// @Router /user/v1/users/{id} [get]
func decodeGetUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getUserRequest{ID: id}, nil
}

// GetUsers godoc
// @Summary Get details of all Users
// @Description Get details of all Users
// @Tags User
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Quantity of items"
// @Success 200 {object} []User
// @Router /user/v1/users [get]
func decodeGetUsersRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page == 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return getUsersRequest{Offset: offset, PageSize: pageSize}, nil
}

// PutUser godoc
// @Summary Change or create User item
// @Description Change or create User item
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param data body user.PostUser true "User struct"
// @Success 200
// @Router /user/v1/users/{id} [put]
func decodePutUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return nil, err
	}
	return putUserRequest{
		ID:   id,
		User: user,
	}, nil
}

// PatchUser godoc
// @Summary Change User item
// @Description Change User item
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param data body user.PostUser true "User struct"
// @Success 200
// @Router /user/v1/users/{id} [patch]
func decodePatchUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return nil, err
	}
	return patchUserRequest{
		ID:   id,
		User: user,
	}, nil
}

// DeleteUser godoc
// @Summary Delete User item
// @Description Delete User item
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200
// @Router /user/v1/users/{id} [delete]
func decodeDeleteUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deleteUserRequest{ID: id}, nil
}

func encodePostUserRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("POST").Path("/users/")
	req.URL.Path = "/user/v1/users"
	return encodeRequest(ctx, req, request)
}

func encodeGetUserRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/users/{id}")
	r := request.(getUserRequest)
	userID := url.QueryEscape(r.ID)
	req.URL.Path = "/user/v1/users/" + userID
	return encodeRequest(ctx, req, request)
}

func encodeGetUsersRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/users")
	req.URL.Path = "/user/v1/users"
	return encodeRequest(ctx, req, request)
}

func encodePutUserRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("PUT").Path("/users/{id}")
	r := request.(putUserRequest)
	userID := url.QueryEscape(r.ID)
	req.URL.Path = "/user/v1/users/" + userID
	return encodeRequest(ctx, req, request)
}

func encodePatchUserRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("PATCH").Path("/users/{id}")
	r := request.(patchUserRequest)
	userID := url.QueryEscape(r.ID)
	req.URL.Path = "/user/v1/users/" + userID
	return encodeRequest(ctx, req, request)
}

func encodeDeleteUserRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("DELETE").Path("/users/{id}")
	r := request.(deleteUserRequest)
	userID := url.QueryEscape(r.ID)
	req.URL.Path = "/user/v1/users/" + userID
	return encodeRequest(ctx, req, request)
}

func decodePostUserResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response postUserResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetUserResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getUserResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetUsersResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getUsersResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePutUserResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response putUserResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePatchUserResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response patchUserResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeDeleteUserResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response deleteUserResponse
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
// usersvc endpoints require mutating the HTTP method and request path.
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
