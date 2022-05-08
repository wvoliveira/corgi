package auth

import (
	"errors"
	"net/http"

	"github.com/elga-io/corgi/internal/app/entity"
)

type logoutRequest struct {
	Token entity.Token
}

func decodeLogout(r *http.Request) (req logoutRequest, err error) {
	ctx := r.Context()

	userID := ctx.Value("user_id")
	if userID == nil {
		return req, errors.New("impossible to get user_id from context")
	}

	refreshTokenID := ctx.Value("refresh_token_id")
	if refreshTokenID == nil {
		return req, errors.New("impossible to get refresh_token_id from context")
	}

	req.Token.ID = refreshTokenID.(string)
	req.Token.UserID = userID.(string)
	return req, nil
}
