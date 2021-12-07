package google

import (
	"net/http"
)

type loginRequest struct{}

func decodeLoginRequest(r *http.Request) (req loginRequest, err error) {
	return req, nil
}

type callbackRequest struct {
	State    string   //state=state
	Code     string   // code=4/0AX4XfWjLX8K0dMCvLgUA86jMy_nTuRhuAKLtxGSC0gFvD8xWiNx-JjEDZ-XX4c93Wq1wzg
	Scopes   []string // scope=email%20profile%20https://www.googleapis.com/auth/userinfo.email%20https://www.googleapis.com/auth/userinfo.profile%20openid
	AuthUser string   // authuser=0
	Domain   string   // hd = elga.io
	Prompt   string   //prompt = consent
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
