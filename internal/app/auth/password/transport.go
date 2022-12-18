package password

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(rg *gin.RouterGroup) {
	r := rg.Group("/auth/password")

	r.POST("/login", s.HTTPLogin)
	r.POST("/register", s.HTTPRegister)
}

func (s service) HTTPLogin(c *gin.Context) {

	dr, err := decodeLogin(c)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	identity := model.Identity{
		Provider: "email",
		UID:      dr.Email,
		Password: dr.Password,
	}

	user, err := s.Login(c, identity)

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
}

func (s service) HTTPRegister(c *gin.Context) {

	dr, err := decodeRegister(c)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	err = s.Register(c, model.Identity{
		Provider: "email",
		UID:      dr.Email,
		Password: dr.Password,
	})

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	encodeRegister(c)
}
