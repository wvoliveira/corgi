package password

import (
	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

type loginResponse struct {
	Name         string `json:"name"`
	Role         string `json:"role"`
	Active       bool   `json:"active"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func encodeLogin(name, role string, active bool, accessToken, refreshToken string) (r loginResponse) {
	r.Name = name
	r.Role = role
	r.Active = active
	r.AccessToken = accessToken
	r.RefreshToken = refreshToken
	return
}

func encodeRegister(c *gin.Context) {
	c.JSON(200, response.Response{Status: "successful"})
}
