package urls

// The URL is just over HTTP, so we just have a single transport.go.

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
	httptransport "github.com/go-kit/kit/transport/http"
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
// Useful in a urlsvc server.
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// GET     /urls/                          list urls
	// POST    /urls/                          adds another url
	// GET     /urls/:id                       retrieves the given url by id
	// PUT     /urls/:id                       post updated url information about the url
	// PATCH   /urls/:id                       partial updated url information
	// DELETE  /urls/:id                       remove the given url

	// GET     /swagger                        swagger specification
	// GET     /                               web ui application

	// GET     /health/ready                   it corresponds to the Kubernetes readiness probe
	// GET     /health/live                    this endpoint corresponds to the Kubernetes liveness probe, which automatically restarts the pod if the check fails.

	r.Methods("GET").Path("/{urls:urls\\/?}").Handler(httptransport.NewServer(
		e.GetURLsEndpoint,
		decodeGetURLsRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/{urls:urls\\/?}").Handler(httptransport.NewServer(
		e.PostURLEndpoint,
		decodePostURLRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/urls/{id}").Handler(httptransport.NewServer(
		e.GetURLEndpoint,
		decodeGetURLRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PUT").Path("/urls/{id}").Handler(httptransport.NewServer(
		e.PutURLEndpoint,
		decodePutURLRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PATCH").Path("/urls/{id}").Handler(httptransport.NewServer(
		e.PatchURLEndpoint,
		decodePatchURLRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("/urls/{id}").Handler(httptransport.NewServer(
		e.DeleteURLEndpoint,
		decodeDeleteURLRequest,
		encodeResponse,
		options...,
	))

	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	// Healthcheck endpoints
	health := healthcheck.NewHandler()
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))

	health.AddReadinessCheck(
		"upstream-dep-dns",
		healthcheck.DNSResolveCheck("localhost", 50*time.Millisecond))

	r.Path("/health/{ready:ready\\/?}").HandlerFunc(health.ReadyEndpoint)
	r.Path("/health/{live:live\\/?}").HandlerFunc(health.LiveEndpoint)

	return r
}

// PostURL godoc
// @Summary Add a new URL
// @Description Add a new URL
// @Tags URLs
// @Accept json
// @Produce json
// @Param data body urls.PostURL true "URL struct"
// @Success 200
// @Router /urls [post]
func decodePostURLRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postURLRequest
	if e := json.NewDecoder(r.Body).Decode(&req.URL); e != nil {
		return nil, e
	}
	return req, nil
}

// GetURL godoc
// @Summary Get an URL
// @Description Get URL by ID
// @Tags URLs
// @Accept json
// @Produce json
// @Param id path string true "URL ID"
// @Success 200
// @Failure 404
// @Router /urls/{id} [get]
func decodeGetURLRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getURLRequest{ID: id}, nil
}

// GetURLs godoc
// @Summary Get details of all URLs
// @Description Get details of all URLs
// @Tags URLs
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Quantity of items"
// @Success 200 {object} []URL
// @Router /urls [get]
func decodeGetURLsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
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
	return getURLsRequest{Offset: offset, PageSize: pageSize}, nil
}

// PutURL godoc
// @Summary Change or create URL item
// @Description Change or create URL item
// @Tags URLs
// @Accept json
// @Produce json
// @Param id path string true "URL ID"
// @Param data body urls.PostURL true "URL struct"
// @Success 200
// @Router /urls/{id} [put]
func decodePutURLRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var url URL
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		return nil, err
	}
	return putURLRequest{
		ID:  id,
		URL: url,
	}, nil
}

// PatchURL godoc
// @Summary Change URL item
// @Description Change URL item
// @Tags URLs
// @Accept json
// @Produce json
// @Param id path string true "URL ID"
// @Param data body urls.PostURL true "URL struct"
// @Success 200
// @Router /urls/{id} [patch]
func decodePatchURLRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var url URL
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		return nil, err
	}
	return patchURLRequest{
		ID:  id,
		URL: url,
	}, nil
}

// DeleteURL godoc
// @Summary Delete URL item
// @Description Delete URL item
// @Tags URLs
// @Accept json
// @Produce json
// @Param id path string true "URL ID"
// @Success 200
// @Router /urls/{id} [delete]
func decodeDeleteURLRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deleteURLRequest{ID: id}, nil
}

func encodePostURLRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("POST").Path("/urls/")
	req.URL.Path = "/urls/"
	return encodeRequest(ctx, req, request)
}

func encodeGetURLRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/urls/{id}")
	r := request.(getURLRequest)
	urlID := url.QueryEscape(r.ID)
	req.URL.Path = "/urls/" + urlID
	return encodeRequest(ctx, req, request)
}

func encodeGetURLsRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/urls")
	req.URL.Path = "/urls"
	return encodeRequest(ctx, req, request)
}

func encodePutURLRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("PUT").Path("/urls/{id}")
	r := request.(putURLRequest)
	urlID := url.QueryEscape(r.ID)
	req.URL.Path = "/urls/" + urlID
	return encodeRequest(ctx, req, request)
}

func encodePatchURLRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("PATCH").Path("/urls/{id}")
	r := request.(patchURLRequest)
	urlID := url.QueryEscape(r.ID)
	req.URL.Path = "/urls/" + urlID
	return encodeRequest(ctx, req, request)
}

func encodeDeleteURLRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("DELETE").Path("/urls/{id}")
	r := request.(deleteURLRequest)
	urlID := url.QueryEscape(r.ID)
	req.URL.Path = "/urls/" + urlID
	return encodeRequest(ctx, req, request)
}

func decodePostURLResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response postURLResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetURLResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getURLResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetURLsResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getURLsResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePutURLResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response putURLResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePatchURLResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response patchURLResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeDeleteURLResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response deleteURLResponse
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
// urlsvc endpoints require mutating the HTTP method and request path.
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
