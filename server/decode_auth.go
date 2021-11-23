package server

import (
	"encoding/json"
	"net/http"
)

func decodeAuthLoginRequest(r *http.Request) (req authLoginRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, err
	}
	return req, nil
}

func decodeAuthRegisterRequest(r *http.Request) (req signUpRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, err
	}
	return req, nil
}
