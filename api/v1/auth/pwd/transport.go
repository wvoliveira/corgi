package pwd

// The "pwd" is just over HTTP, so we just have a single transport.go.

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
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
// Useful in a "pwd" service.
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	signInPwdHandler := kithttp.NewServer(
		e.SignInPwdEndpoint,
		decodeSignInPwdRequest,
		encodeSignInPwdResponse,
		options...,
	)

	signUpPwdHandler := kithttp.NewServer(
		e.SignUpPwdEndpoint,
		decodeSignUpPwdRequest,
		encodeResponse,
		options...,
	)

	health := healthcheck.NewHandler()
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))
	health.AddReadinessCheck("upstream-dep-dns", healthcheck.DNSResolveCheck("localhost", 50*time.Millisecond))

	r.Handle("/auth/pwd/v1/signin", signInPwdHandler).Methods("POST")
	r.Handle("/auth/pwd/v1/signup", signUpPwdHandler).Methods("POST")

	r.HandleFunc("/auth/pwd/v1/health/ready", health.ReadyEndpoint)
	r.HandleFunc("/auth/pwd/v1/health/live", health.LiveEndpoint)
	return r
}

// SignInPwd godoc
// @Summary Authenticate the user
// @Description Authenticate the user
// @Tags Pwd
// @Accept json
// @Produce json
// @Param data body pwd.signInPwdRequest true "Pwd struct"
// @Success 200
// @Router /auth/pwd/v1/signin [post]
func decodeSignInPwdRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req signInPwdRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

// SignUpPwd godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags Pwd
// @Accept json
// @Produce json
// @Param data body pwd.signInPwdRequest true "Pwd struct"
// @Success 200
// @Router /auth/pwd/v1/signup [post]
func decodeSignUpPwdRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req signUpPwdRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

func encodeSignInPwdResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}

	var sessionToken string

	// Get session_token
	e, ok := response.(signInPwdResponse)
	if ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	sessionToken = e.SessionToken

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds, the same as the cache
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(300 * time.Second),
	})

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
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
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case io.EOF:
		return http.StatusBadRequest
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists, ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
