package user

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type findRequest struct {
	UserID string
}

type updateRequest struct {
	UserID string `json:"-"`
	Name   string `json:"name"`
}

func decodeFind(c *gin.Context) (req findRequest, err error) {
	userID, _ := c.Get("user_id")
	req.UserID = userID.(string)
	return req, nil
}

func decodeUpdate(c *gin.Context) (req updateRequest, err error) {
	userID, _ := c.Get("user_id")
	if err = json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return req, err
	}
	req.UserID = userID.(string)
	return req, nil
}
