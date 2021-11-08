package server

// The "account" is just over HTTP, so we just have a single transport.go.

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"

	"github.com/heptiolabs/healthcheck"
)

/*
MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
Useful in a "account" service.
*/
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerBefore(putRequestInCtx),
	}

	/*
		Auth handlers.
	*/
	signInHandler := kithttp.NewServer(
		e.SignInEndpoint,
		decodeSignInRequest,
		encodeSignInResponse,
		options...,
	)

	signUpHandler := kithttp.NewServer(
		e.SignUpEndpoint,
		decodeSignUpRequest,
		encodeResponse,
		options...,
	)

	/*
		Account handlers.
	*/

	addAccountHandler := kithttp.NewServer(
		AuthMiddleware()(e.AddAccountEndpoint),
		decodeAddAccountRequest,
		encodeResponse,
		options...,
	)

	findAccountByIDHandler := kithttp.NewServer(
		AuthMiddleware()(e.FindAccountByIDEndpoint),
		decodeFindAccountByIDRequest,
		encodeResponse,
		options...,
	)

	findAccountsHandler := kithttp.NewServer(
		AuthMiddleware()(e.FindAccountsEndpoint),
		decodeFindAccountsRequest,
		encodeResponse,
		options...,
	)

	updateOrCreateAccountHandler := kithttp.NewServer(
		AuthMiddleware()(e.UpdateOrCreateAccountEndpoint),
		decodeUpdateOrCreateAccountRequest,
		encodeResponse,
		options...,
	)

	updateAccountHandler := kithttp.NewServer(
		AuthMiddleware()(e.UpdateAccountEndpoint),
		decodeUpdateAccountRequest,
		encodeResponse,
		options...,
	)

	deleteAccountHandler := kithttp.NewServer(
		AuthMiddleware()(e.DeleteAccountEndpoint),
		decodeDeleteAccountRequest,
		encodeResponse,
		options...,
	)

	/*
		URL handlers.
	*/

	addURLHandler := kithttp.NewServer(
		AuthMiddleware()(e.AddURLEndpoint),
		decodeAddURLRequest,
		encodeResponse,
		options...,
	)

	findURLByIDHandler := kithttp.NewServer(
		AuthMiddleware()(e.FindURLByIDEndpoint),
		decodeFindURLByIDRequest,
		encodeResponse,
		options...,
	)

	findURLsHandler := kithttp.NewServer(
		AuthMiddleware()(e.FindURLsEndpoint),
		decodeFindURLsRequest,
		encodeResponse,
		options...,
	)

	updateOrCreateURLHandler := kithttp.NewServer(
		AuthMiddleware()(e.UpdateOrCreateURLEndpoint),
		decodeUpdateOrCreateURLRequest,
		encodeResponse,
		options...,
	)

	updateURLHandler := kithttp.NewServer(
		AuthMiddleware()(e.UpdateURLEndpoint),
		decodeUpdateURLRequest,
		encodeResponse,
		options...,
	)

	deleteURLHandler := kithttp.NewServer(
		AuthMiddleware()(e.DeleteURLEndpoint),
		decodeDeleteURLRequest,
		encodeResponse,
		options...,
	)

	/*
		Health check functions and endpoints.
		TODO: separate by service.
	*/
	health := healthcheck.NewHandler()
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))
	health.AddReadinessCheck("upstream-dep-dns", healthcheck.DNSResolveCheck("localhost", 50*time.Millisecond))

	r.HandleFunc("/health/ready", health.ReadyEndpoint)
	r.HandleFunc("/health/live", health.LiveEndpoint)

	/*
		Auth with password endpoints and methods.
	*/
	r.Handle("/api/v1/signin", signInHandler).Methods("POST")
	r.Handle("/api/v1/signup", signUpHandler).Methods("POST")

	/*
		Account endpoints and methods.
	*/
	r.Handle("/api/v1/accounts", addAccountHandler).Methods("POST")
	r.Handle("/api/v1/accounts/{id}", findAccountByIDHandler).Methods("GET")
	r.Handle("/api/v1/accounts", findAccountsHandler).Methods("GET")
	r.Handle("/api/v1/accounts/{id}", updateOrCreateAccountHandler).Methods("PUT")
	r.Handle("/api/v1/accounts/{id}", updateAccountHandler).Methods("PATCH")
	r.Handle("/api/v1/accounts/{id}", deleteAccountHandler).Methods("DELETE")

	/*
		URL endpoints and methods.
	*/
	r.Handle("/api/v1/urls", addURLHandler).Methods("POST")
	r.Handle("/api/v1/urls/{id}", findURLByIDHandler).Methods("GET")
	r.Handle("/api/v1/urls", findURLsHandler).Methods("GET")
	r.Handle("/api/v1/urls/{id}", updateOrCreateURLHandler).Methods("PUT")
	r.Handle("/api/v1/urls/{id}", updateURLHandler).Methods("PATCH")
	r.Handle("/api/v1/urls/{id}", deleteURLHandler).Methods("DELETE")

	return r
}
