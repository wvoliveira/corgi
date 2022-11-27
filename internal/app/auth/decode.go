package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
)

type logoutRequest struct {
	Token entity.Token
}

func decodeLogout(c *gin.Context) (req logoutRequest, err error) {
	// TODO: check this c.Value(...)
	data := c.Value(entity.IdentityInfo{})
	if data == nil {
		return req, errors.New("impossible to get identity from context")
	}

	ii := data.(entity.IdentityInfo)

	req.Token.ID = ii.RefreshTokenID
	req.Token.UserID = ii.UserID
	return req, nil
}
