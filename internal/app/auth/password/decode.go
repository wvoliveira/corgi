package password

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func decodeLogin(c *gin.Context) (req loginRequest, err error) {
	if err = json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return req, err
	}
	return req, nil
}

func decodeRegister(c *gin.Context) (req registerRequest, err error) {
	err = json.NewDecoder(c.Request.Body).Decode(&req)
	return req, err
}
