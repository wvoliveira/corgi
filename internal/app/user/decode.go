package user

import (
	"encoding/json"
	"errors"
	"net/http"
)

type findRequest struct {
	UserID string
}

type updateRequest struct {
	UserID string `json:"-"`
	Name   string `json:"name"`
}

func decodeFind(r *http.Request) (req findRequest, err error) {
	ctx := r.Context()

	userID := ctx.Value("user_id")
	if userID == nil {
		return req, errors.New("impossible to get user_id from context")
	}

	req.UserID = userID.(string)
	return req, nil
}

func decodeUpdate(r *http.Request) (req updateRequest, err error) {
	ctx := r.Context()

	userID := ctx.Value("user_id")
	if userID == nil {
		return req, errors.New("impossible to get user_id from context")
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, err
	}
	req.UserID = userID.(string)
	return req, nil
}
