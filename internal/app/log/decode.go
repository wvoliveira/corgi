package log

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

type addRequest struct {
	user model.User `json:"-"`
	log  any
}

func decodeAdd(c *gin.Context) (user model.User, payload any, err error) {
	data, exists := c.Get("user")

	raw, err := c.GetRawData()
	if err != nil {
		return user, payload, err
	}

	print("RAW:")
	print(string(raw))

	if !exists {
		return user, payload, errors.New("impossible to get user from context")
	}

	user = data.(model.User)

	if err = json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
		return user, payload, err
	}

	print(payload)

	return
}
