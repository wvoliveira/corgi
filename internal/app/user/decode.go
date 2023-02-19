package user

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
)

type findMeRequest struct {
	whoID string
}

type updateMeRequest struct {
	whoID string
	User  struct {
		Name string `json:"name"`
	}
}

type findByIDorUsernameRequest struct {
	whoID        string
	IDorUsername string `uri:"id_username" binding:"required"`
}

type updateIDorUsernameRequest struct {
	whoID        string
	IDorUsername string `uri:"id_username" binding:"required"`
	User         struct {
		Name string `json:"name"`
	}
}

func decodeFindMe(c *gin.Context) (r findMeRequest, err error) {
	v, ok := c.Get("user_id")
	if !ok {
		return r, errors.New("impossible to know who you are")
	}

	r.whoID = v.(string)
	return
}

func decodeUpdateMe(c *gin.Context) (r updateMeRequest, err error) {
	v, ok := c.Get("user_id")
	if !ok {
		return r, errors.New("impossible to know who you are")
	}

	r.whoID = v.(string)

	err = c.ShouldBindUri(&r)
	if err != nil {
		return r, errors.New("impossible to get user ID or username from URI")
	}

	if err = json.NewDecoder(c.Request.Body).Decode(&r.User); err != nil {
		return r, err
	}
	return
}

func decodeFindByIDorUsername(c *gin.Context) (r findByIDorUsernameRequest, err error) {
	v, ok := c.Get("user_id")
	if !ok {
		return r, errors.New("impossible to know who you are")
	}

	r.whoID = v.(string)

	err = c.ShouldBindUri(&r)
	if err != nil {
		return r, errors.New("impossible to get user ID or username from URI")
	}
	return
}

func decodeUpdateByIDorUsername(c *gin.Context) (r updateIDorUsernameRequest, err error) {
	v, ok := c.Get("user_id")
	if !ok {
		return r, errors.New("impossible to know who you are")
	}

	r.whoID = v.(string)

	err = c.ShouldBindUri(&r)
	if err != nil {
		return r, errors.New("impossible to get user ID or username from URI")
	}

	if err = json.NewDecoder(c.Request.Body).Decode(&r.User); err != nil {
		return r, err
	}
	return
}
