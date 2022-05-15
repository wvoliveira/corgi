package facebook

import (
	"net/http"
)

type loginRequest struct{}

func decodeLoginRequest(r *http.Request) (req loginRequest, err error) {
	return req, nil
}

type callbackRequest struct {
	State    string
	Code     string
	Scopes   []string
	AuthUser string
	Domain   string
	Prompt   string
}

func decodeCallbackRequest(r *http.Request) (req callbackRequest, err error) {
	q := r.URL.Query()
	var scopes []string
	req = callbackRequest{
		q.Get("state"),
		q.Get("code"),
		append(scopes, q.Get("scopes")),
		q.Get("authuser"),
		q.Get("hd"),
		q.Get("prompt"),
	}
	return
}
