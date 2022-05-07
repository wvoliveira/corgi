package token

import (
	"encoding/json"
	"net/http"
)

type refreshRequest struct {
	RefreshTokenID string `json:"refresh_token_id"`
}

func decodeRefreshRequest(r *http.Request) (req refreshRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, err
	}
	return req, nil
}
