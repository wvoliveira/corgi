package password

import (
	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

type loginResponse struct {
	User struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Role     string `json:"role"`
		Active   bool   `json:"active"`
	} `json:"user"`
	Tokens struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	} `json:"tokens"`
}

func encodeLogin(username, name, role string, active bool, accessToken, refreshToken string) (r loginResponse) {
	r.User.Username = username
	r.User.Name = name
	r.User.Role = role
	r.User.Active = active
	r.Tokens.AccessToken = accessToken
	r.Tokens.RefreshToken = refreshToken
	return
}

func encodeRegister(c *gin.Context) {
	c.JSON(200, response.Response{Status: "successful"})
}
