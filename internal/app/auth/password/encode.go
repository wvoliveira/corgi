package password

import (
	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

type loginResponse struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Active   bool   `json:"active"`
	Tokens   struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	} `json:"tokens"`
}

func encodeLogin(username, name, role string, active bool, accessToken, refreshToken string) (r loginResponse) {
	r.Username = username
	r.Name = name
	r.Role = role
	r.Active = active
	r.Tokens.AccessToken = accessToken
	r.Tokens.RefreshToken = refreshToken
	return
}

func encodeRegister(c *gin.Context) {
	c.JSON(200, response.Response{Status: "successful"})
}
