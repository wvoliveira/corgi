package server

import (
	"encoding/json"
	"net/http"
)

func decodeTokenRefreshRequest(r *http.Request) (acc Account, req TokenRefreshRequest, err error) {
	acc = getAccountFromHeaders(r)
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return acc, req, err
	}
	return acc, req, nil
}
