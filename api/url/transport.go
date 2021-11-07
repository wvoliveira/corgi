package url

// The URL is just over HTTP, so we just have a single transport.go.

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"

	"github.com/heptiolabs/healthcheck"
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
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	getURLsHandler := kithttp.NewServer(
		e.GetURLsEndpoint,
		decodeGetURLsRequest,
		encodeResponse,
		options...,
	)

	postURLHandler := kithttp.NewServer(
		e.PostURLEndpoint,
		decodePostURLRequest,
		encodeResponse,
		options...,
	)

	getURLHandler := kithttp.NewServer(
		e.GetURLEndpoint,
		decodeGetURLRequest,
		encodeResponse,
		options...,
	)

	putURLHandler := kithttp.NewServer(
		e.PutURLEndpoint,
		decodePutURLRequest,
		encodeResponse,
		options...,
	)

	patchURLHandler := kithttp.NewServer(
		e.PatchURLEndpoint,
		decodePatchURLRequest,
		encodeResponse,
		options...,
	)

	deleteURLHandler := kithttp.NewServer(
		e.DeleteURLEndpoint,
		decodeDeleteURLRequest,
		encodeResponse,
		options...,
	)

	health := healthcheck.NewHandler()
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))
	health.AddReadinessCheck("upstream-dep-dns", healthcheck.DNSResolveCheck("localhost", 50*time.Millisecond))

	r.Handle("/url/v1/urls", getURLsHandler).Methods("GET")
	r.Handle("/url/v1/urls", postURLHandler).Methods("POST")
	r.Handle("/url/v1/urls/{id}", getURLHandler).Methods("GET")
	r.Handle("/url/v1/urls/{id}", putURLHandler).Methods("PUT")
	r.Handle("/url/v1/urls/{id}", patchURLHandler).Methods("PATCH")
	r.Handle("/url/v1/urls/{id}", deleteURLHandler).Methods("DELETE")

	r.HandleFunc("/url/v1/health/ready", health.ReadyEndpoint)
	r.HandleFunc("/url/v1/health/live", health.LiveEndpoint)

	return r
}

// PostURL godoc
// @Summary Add a new URL
// @Description Add a new URL
// @Tags URL
// @Accept json
// @Produce json
// @Param data body url.postURLRequest true "URL struct"
// @Success 200
// @Router /url/v1/urls [post]
func decodePostURLRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req postURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req.URL); err != nil {
		return nil, err
	}
	return req, nil
}

// GetURL godoc
// @Summary Get an URL
// @Description Get URL by ID
// @Tags URL
// @Accept json
// @Produce json
// @Param id path string true "URL ID"
// @Success 200
// @Failure 404
// @Router /url/v1/urls/{id} [get]
func decodeGetURLRequest(_ context.Context, r *http.Request) (interface{}, error) {
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
// @Tags URL
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Quantity of items"
// @Success 200 {object} []URL
// @Router /url/v1/urls [get]
func decodeGetURLsRequest(_ context.Context, r *http.Request) (interface{}, error) {
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
// @Tags URL
// @Accept json
// @Produce json
// @Param id path string true "URL ID"
// @Param data body url.postURLRequest true "URL struct"
// @Success 200
// @Router /url/v1/urls/{id} [put]
func decodePutURLRequest(_ context.Context, r *http.Request) (interface{}, error) {
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
// @Tags URL
// @Accept json
// @Produce json
// @Param id path string true "URL ID"
// @Param data body url.postURLRequest true "URL struct"
// @Success 200
// @Router /url/v1/urls/{id} [patch]
func decodePatchURLRequest(_ context.Context, r *http.Request) (interface{}, error) {
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
// @Tags URL
// @Accept json
// @Produce json
// @Param id path string true "URL ID"
// @Success 200
// @Router /url/v1/urls/{id} [delete]
func decodeDeleteURLRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deleteURLRequest{ID: id}, nil
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
