package pwd

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints collects all of the endpoints that compose a Pwd service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
//
// In a server, it's useful for functions that need to operate on a per-endpoint
// basis. For example, you might pass an Endpoints to a function that produces
// an http.Handler, with each method (endpoint) wired up to a specific path. (It
// is probably a mistake in design to invoke the Service methods on the
// Endpoints struct in a server.)
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them
// into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	SignInPwdEndpoint endpoint.Endpoint
	SignUpPwdEndpoint endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service. Useful in a Pwdsvc
// server.
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		SignInPwdEndpoint: MakeSignInPwdEndpoint(s),
		SignUpPwdEndpoint: MakeSignUpPwdEndpoint(s),
	}
}

// SignInPwd implements Service. Primarily useful in a client.
func (e Endpoints) SignInPwd(ctx context.Context, p Pwd) error {
	request := signInPwdRequest{Pwd: p}
	response, err := e.SignInPwdEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(signInPwdResponse)
	return resp.Err
}

// SignUpPwd implements Service. Primarily useful in a client.
func (e Endpoints) SignUpPwd(ctx context.Context, p Pwd) error {
	request := signUpPwdRequest{Pwd: p}
	response, err := e.SignUpPwdEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(signUpPwdResponse)
	return resp.Err
}

// MakeSignInPwdEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeSignInPwdEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(signInPwdRequest)
		p, e := s.SignInPwd(ctx, req.Pwd)
		return signInPwdResponse{SessionToken: p.SessionToken, Err: e}, nil
	}
}

// MakeSignUpPwdEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeSignUpPwdEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(signUpPwdRequest)
		e := s.SignUpPwd(ctx, req.Pwd)
		return signUpPwdResponse{Err: e}, nil
	}
}

// We have two options to return errors from the business logic.
//
// We could return the error via the endpoint itself. That makes certain things
// a little bit easier, like providing non-200 HTTP responses to the client. But
// Go kit assumes that endpoint errors are (or may be treated as)
// transport-domain errors. For example, an endpoint error will count against a
// circuit breaker error count.
//
// Therefore, it's often better to return service (business logic) errors in the
// response object. This means we have to do a bit more work in the HTTP
// response encoder to detect e.g. a not-found error and provide a proper HTTP
// status code. That work is done with the errorer interface, in transport.go.
// Response types that may contain business-logic errors implement that
// interface.

type signInPwdRequest struct {
	Pwd
}

type signUpPwdRequest struct {
	Pwd
}

type signInPwdResponse struct {
	SessionToken string `json:"-"`
	Err          error  `json:"err,omitempty"`
}

type signUpPwdResponse struct {
	Err error `json:"err,omitempty"`
}

func (r signInPwdResponse) error() error { return r.Err }

func (r signUpPwdResponse) error() error { return r.Err }
