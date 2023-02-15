package user

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

func decodeFind(c *gin.Context) (user model.User, err error) {
	v, ok := c.Get("user_id")

	if !ok {
		return user, errors.New("impossible to get user from context")
	}

	user.ID = v.(string)
	return
}

func decodeUpdate(c *gin.Context) (user model.User, err error) {
	v, ok := c.Get("user_id")

	if !ok {
		return user, errors.New("impossible to get user from context")
	}

	user.ID = v.(string)

	if err = json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		return user, err
	}

	return
}
