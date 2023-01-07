package facebook

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(rg *gin.RouterGroup) {
	r := rg.Group("/auth/facebook")

	r.GET("/login", s.HTTPLogin)
	r.GET("/callback", s.HTTPCallback)
}

func (s service) HTTPLogin(c *gin.Context) {

	// Remove this and set in decode function
	accessToken, _ := c.GetQuery("access_token")

	schema := "http"
	if c.Request.TLS != nil {
		schema = "https"
	}

	// Ex.: http://localhost:8081/api/auth/facebook/callback
	callbackURL := fmt.Sprintf("%s://%s", schema, c.Request.Host+"/auth/facebook/callback")

	user, redirectURL, err := s.Login(c, accessToken, callbackURL)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	if user.ID != "" {
		session := sessions.Default(c)
		session.Set("user", user)

		err = session.Save()

		if err != nil {
			e.EncodeError(c, err)
			return
		}

		c.JSON(200, response.Response{
			Status: "successful",
			Data:   user,
		})

		return
	}

	c.Redirect(http.StatusMovedPermanently, redirectURL)
}

func (s service) HTTPCallback(c *gin.Context) {

	dr, err := decodeCallbackRequest(c)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	schema := "http"
	if c.Request.TLS != nil {
		schema = "https"
	}

	callbackURL := fmt.Sprintf("%s://%s", schema, c.Request.Host+"/auth/facebook/callback")

	user, err := s.Callback(c, callbackURL, dr)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	session := sessions.Default(c)
	session.Set("user", user)

	err = session.Save()

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	data := gin.H{
		"name":   user.Name,
		"role":   user.Role,
		"active": user.Active,
	}

	c.JSON(200, response.Response{
		Status: "successful",
		Data:   data,
	})

	c.Redirect(http.StatusMovedPermanently, viper.GetString("app.redirect_url"))
}
