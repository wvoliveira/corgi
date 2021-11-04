package profile

// The "profile" is just over HTTP, so we just have a single transport.go.

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
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
// Useful in a "profile" service.
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	getProfilesHandler := kithttp.NewServer(
		e.GetProfilesEndpoint,
		decodeGetProfilesRequest,
		encodeResponse,
		options...,
	)

	postProfileHandler := kithttp.NewServer(
		e.PostProfileEndpoint,
		decodePostProfileRequest,
		encodeResponse,
		options...,
	)

	getProfileHandler := kithttp.NewServer(
		e.GetProfileEndpoint,
		decodeGetProfileRequest,
		encodeResponse,
		options...,
	)

	putProfileHandler := kithttp.NewServer(
		e.PutProfileEndpoint,
		decodePutProfileRequest,
		encodeResponse,
		options...,
	)

	patchProfileHandler := kithttp.NewServer(
		e.PatchProfileEndpoint,
		decodePatchProfileRequest,
		encodeResponse,
		options...,
	)

	deleteProfileHandler := kithttp.NewServer(
		e.DeleteProfileEndpoint,
		decodeDeleteProfileRequest,
		encodeResponse,
		options...,
	)

	health := healthcheck.NewHandler()
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))
	health.AddReadinessCheck("upstream-dep-dns", healthcheck.DNSResolveCheck("localhost", 50*time.Millisecond))

	r.Handle("/profile/v1/profiles", getProfilesHandler).Methods("GET")
	r.Handle("/profile/v1/profiles", postProfileHandler).Methods("POST")
	r.Handle("/profile/v1/profiles/{id}", getProfileHandler).Methods("GET")
	r.Handle("/profile/v1/profiles/{id}", putProfileHandler).Methods("PUT")
	r.Handle("/profile/v1/profiles/{id}", patchProfileHandler).Methods("PATCH")
	r.Handle("/profile/v1/profiles/{id}", deleteProfileHandler).Methods("DELETE")

	r.HandleFunc("/profile/v1/health/ready", health.ReadyEndpoint)
	r.HandleFunc("/profile/v1/health/live", health.LiveEndpoint)

	r.PathPrefix("/profile/v1/swagger").Handler(httpSwagger.WrapHandler)

	return r
}

// PostProfile godoc
// @Summary Add a new Profile
// @Description Add a new Profile
// @Tags Profile
// @Accept json
// @Produce json
// @Param data body profile.postProfileRequest true "Profile struct"
// @Success 200
// @Router /profile/v1/profiles [post]
func decodePostProfileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postProfileRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Profile); e != nil {
		return nil, e
	}
	return req, nil
}

// GetProfile godoc
// @Summary Get an Profile
// @Description Get Profile by ID
// @Tags Profile
// @Accept json
// @Produce json
// @Param id path string true "Profile ID"
// @Success 200
// @Failure 404
// @Router /profile/v1/profiles/{id} [get]
func decodeGetProfileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getProfileRequest{ID: id}, nil
}

// GetProfiles godoc
// @Summary Get details of all Profiles
// @Description Get details of all Profiles
// @Tags Profile
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Quantity of items"
// @Success 200 {object} []Profile
// @Router /profile/v1/profiles [get]
func decodeGetProfilesRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
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
	return getProfilesRequest{Offset: offset, PageSize: pageSize}, nil
}

// PutProfile godoc
// @Summary Change or create Profile item
// @Description Change or create Profile item
// @Tags Profile
// @Accept json
// @Produce json
// @Param id path string true "Profile ID"
// @Param data body profile.postProfileRequest true "Profile struct"
// @Success 200
// @Router /profile/v1/profiles/{id} [put]
func decodePutProfileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var profile Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		return nil, err
	}
	return putProfileRequest{
		ID:      id,
		Profile: profile,
	}, nil
}

// PatchProfile godoc
// @Summary Change Profile item
// @Description Change Profile item
// @Tags Profile
// @Accept json
// @Produce json
// @Param id path string true "Profile ID"
// @Param data body profile.postProfileRequest true "Profile struct"
// @Success 200
// @Router /profile/v1/profiles/{id} [patch]
func decodePatchProfileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var profile Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		return nil, err
	}
	return patchProfileRequest{
		ID:      id,
		Profile: profile,
	}, nil
}

// DeleteProfile godoc
// @Summary Delete Profile item
// @Description Delete Profile item
// @Tags Profile
// @Accept json
// @Produce json
// @Param id path string true "Profile ID"
// @Success 200
// @Router /profile/v1/profiles/{id} [delete]
func decodeDeleteProfileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deleteProfileRequest{ID: id}, nil
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
