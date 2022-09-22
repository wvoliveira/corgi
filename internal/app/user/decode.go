package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/wvoliveira/corgi/internal/pkg/entity"
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

	data := ctx.Value(entity.IdentityInfo{})
	if data == nil {
		return req, errors.New("impossible to get identity from context")
	}

	ii := data.(entity.IdentityInfo)

	req.UserID = ii.UserID
	return req, nil
}

func decodeUpdate(r *http.Request) (req updateRequest, err error) {
	ctx := r.Context()

	data := ctx.Value(entity.IdentityInfo{})
	if data == nil {
		return req, errors.New("impossible to get identity from context")
	}

	ii := data.(entity.IdentityInfo)

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, err
	}
	req.UserID = ii.UserID
	return req, nil
}
