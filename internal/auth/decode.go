package auth

import (
	"errors"
	"github.com/elga-io/corgi/internal/entity"
	"github.com/gin-gonic/gin"
)

type logoutRequest struct {
	Token entity.Token
}

func decodeLogout(c *gin.Context) (req logoutRequest, err error) {
	userID, ok := c.Get("user_id")
	if !ok {
		return req, errors.New("impossible to get user_id from gin context")
	}
	req.Token.UserID = userID.(string)
	return req, nil
}
