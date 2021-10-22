package user

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Endpoints collects all of the endpoints that compose a User service. It's
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
	PostUserEndpoint   endpoint.Endpoint
	GetUserEndpoint    endpoint.Endpoint
	GetUsersEndpoint   endpoint.Endpoint
	PutUserEndpoint    endpoint.Endpoint
	PatchUserEndpoint  endpoint.Endpoint
	DeleteUserEndpoint endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service. Useful in a Usersvc
// server.
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PostUserEndpoint:   MakePostUserEndpoint(s),
		GetUserEndpoint:    MakeGetUserEndpoint(s),
		GetUsersEndpoint:   MakeGetUsersEndpoint(s),
		PutUserEndpoint:    MakePutUserEndpoint(s),
		PatchUserEndpoint:  MakePatchUserEndpoint(s),
		DeleteUserEndpoint: MakeDeleteUserEndpoint(s),
	}
}

// MakeClientEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the remote instance, via a transport/http.Client.
// Useful in a URLsvc client.
func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	// Note that the request encoders need to modify the request URL, changing
	// the path. That's fine: we simply need to provide specific encoders for
	// each endpoint.

	return Endpoints{
		PostURLEndpoint:   httptransport.NewClient("POST", tgt, encodePostURLRequest, decodePostURLResponse, options...).Endpoint(),
		GetURLEndpoint:    httptransport.NewClient("GET", tgt, encodeGetURLRequest, decodeGetURLResponse, options...).Endpoint(),
		GetURLsEndpoint:   httptransport.NewClient("GET", tgt, encodeGetURLsRequest, decodeGetURLsResponse, options...).Endpoint(),
		PutURLEndpoint:    httptransport.NewClient("PUT", tgt, encodePutURLRequest, decodePutURLResponse, options...).Endpoint(),
		PatchURLEndpoint:  httptransport.NewClient("PATCH", tgt, encodePatchURLRequest, decodePatchURLResponse, options...).Endpoint(),
		DeleteURLEndpoint: httptransport.NewClient("DELETE", tgt, encodeDeleteURLRequest, decodeDeleteURLResponse, options...).Endpoint(),
	}, nil
}

// PostURL implements Service. Primarily useful in a client.
func (e Endpoints) PostURL(ctx context.Context, p URL) error {
	request := postURLRequest{URL: p}
	response, err := e.PostURLEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(postURLResponse)
	return resp.Err
}

// GetURL implements Service. Primarily useful in a client.
func (e Endpoints) GetURL(ctx context.Context, id string) (URL, error) {
	request := getURLRequest{ID: id}
	response, err := e.GetURLEndpoint(ctx, request)
	if err != nil {
		return URL{}, err
	}
	resp := response.(getURLResponse)
	return resp.URL, resp.Err
}

// GetURLs implements Service. Primarily useful in a client.
func (e Endpoints) GetURLs(ctx context.Context, offset, pageSize int) ([]URL, error) {
	request := getURLsRequest{Offset: offset, PageSize: pageSize}
	response, err := e.GetURLsEndpoint(ctx, request)
	if err != nil {
		return []URL{}, err
	}
	resp := response.(getURLsResponse)
	return resp.URLs, resp.Err
}

// PutURL implements Service. Primarily useful in a client.
func (e Endpoints) PutURL(ctx context.Context, id string, p URL) error {
	request := putURLRequest{ID: id, URL: p}
	response, err := e.PutURLEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(putURLResponse)
	return resp.Err
}

// PatchURL implements Service. Primarily useful in a client.
func (e Endpoints) PatchURL(ctx context.Context, id string, p URL) error {
	request := patchURLRequest{ID: id, URL: p}
	response, err := e.PatchURLEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(patchURLResponse)
	return resp.Err
}

// DeleteURL implements Service. Primarily useful in a client.
func (e Endpoints) DeleteURL(ctx context.Context, id string) error {
	request := deleteURLRequest{ID: id}
	response, err := e.DeleteURLEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(deleteURLResponse)
	return resp.Err
}

// MakePostURLEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePostURLEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postURLRequest)
		e := s.PostURL(ctx, req.URL)
		return postURLResponse{Err: e}, nil
	}
}

// MakeGetURLEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetURLEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getURLRequest)
		p, e := s.GetURL(ctx, req.ID)
		return getURLResponse{URL: p, Err: e}, nil
	}
}

// MakeGetURLsEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetURLsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getURLsRequest)
		p, e := s.GetURLs(ctx, req.Offset, req.PageSize)
		return getURLsResponse{URLs: p, Err: e}, nil
	}
}

// MakePutURLEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePutURLEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(putURLRequest)
		e := s.PutURL(ctx, req.ID, req.URL)
		return putURLResponse{Err: e}, nil
	}
}

// MakePatchURLEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePatchURLEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(patchURLRequest)
		e := s.PatchURL(ctx, req.ID, req.URL)
		return patchURLResponse{Err: e}, nil
	}
}

// MakeDeleteURLEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeDeleteURLEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteURLRequest)
		e := s.DeleteURL(ctx, req.ID)
		return deleteURLResponse{Err: e}, nil
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

type postURLRequest struct {
	URL URL
}

type postURLResponse struct {
	Err error `json:"err,omitempty"`
}

func (r postURLResponse) error() error { return r.Err }

type getURLRequest struct {
	ID string
}

type getURLsRequest struct {
	Offset   int
	PageSize int
}

type getURLResponse struct {
	URL URL   `json:"data,omitempty"`
	Err error `json:"error,omitempty"`
}

type getURLsResponse struct {
	URLs []URL `json:"data,omitempty"`
	Err  error `json:"error,omitempty"`
}

func (r getURLResponse) error() error { return r.Err }

type putURLRequest struct {
	ID  string
	URL URL
}

type putURLResponse struct {
	Err error `json:"err,omitempty"`
}

func (r putURLResponse) error() error { return nil }

type patchURLRequest struct {
	ID  string
	URL URL
}

type patchURLResponse struct {
	Err error `json:"err,omitempty"`
}

func (r patchURLResponse) error() error { return r.Err }

type deleteURLRequest struct {
	ID string
}

type deleteURLResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteURLResponse) error() error { return r.Err }
