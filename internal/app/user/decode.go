package user

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
)

func decodeFind(c *gin.Context) (user entity.User, err error) {
	session := sessions.Default(c)
	v := session.Get("user")

	fmt.Println(v)

	if v == nil {
		return user, errors.New("impossible to get user from session")
	}

	err = json.Unmarshal([]byte(v.(string)), &user)
	return
}

func decodeUpdate(c *gin.Context) (user entity.User, err error) {
	data, exists := c.Get("user")

	if !exists {
		return user, errors.New("impossible to get user from context")
	}

	user = data.(entity.User)

	if err = json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		return user, err
	}

	return
}
