package auth

import (
	"errors"
	"net/http"

	"github.com/wvoliveira/corgi/internal/app/entity"
)

type logoutRequest struct {
	Token entity.Token
}

func decodeLogout(r *http.Request) (req logoutRequest, err error) {
	ctx := r.Context()

	data := ctx.Value(entity.IdentityInfo{})
	if data == nil {
		return req, errors.New("impossible to get identity from context")
	}

	ii := data.(entity.IdentityInfo)

	req.Token.ID = ii.RefreshTokenID
	req.Token.UserID = ii.UserID
	return req, nil
}
