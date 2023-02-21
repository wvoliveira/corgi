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

type findByIDRequest struct {
	whoID string
	ID    string `uri:"id" binding:"required"`
}

type updateIDRequest struct {
	whoID string
	ID    string `uri:"id" binding:"required"`
	User  struct {
		Name string `json:"name"`
	}
}

type findByUsernameRequest struct {
	whoID    string
	Username string `uri:"username" binding:"required"`
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
		return r, errors.New("impossible to get user ID from URI")
	}

	if err = json.NewDecoder(c.Request.Body).Decode(&r.User); err != nil {
		return r, err
	}
	return
}

func decodeFindByID(c *gin.Context) (r findByIDRequest, err error) {
	v, ok := c.Get("user_id")
	if !ok {
		return r, errors.New("impossible to know who you are")
	}

	r.whoID = v.(string)

	err = c.ShouldBindUri(&r)
	if err != nil {
		return r, errors.New("impossible to get username from URI")
	}
	return
}

func decodeUpdateByID(c *gin.Context) (r updateIDRequest, err error) {
	v, ok := c.Get("user_id")
	if !ok {
		return r, errors.New("impossible to know who you are")
	}

	r.whoID = v.(string)

	err = c.ShouldBindUri(&r)
	if err != nil {
		return r, errors.New("impossible to get username from URI")
	}

	if err = json.NewDecoder(c.Request.Body).Decode(&r.User); err != nil {
		return r, err
	}
	return
}

func decodeFindByUsername(c *gin.Context) (r findByUsernameRequest, err error) {
	v, ok := c.Get("user_id")
	if !ok {
		return r, errors.New("impossible to know who you are")
	}

	r.whoID = v.(string)

	err = c.ShouldBindUri(&r)
	if err != nil {
		return r, errors.New("impossible to get user ID from URI")
	}
	return
}
