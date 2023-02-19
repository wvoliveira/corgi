package password

import (
	"net/http"

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

	accessToken, refreshToken, user, err := s.Login(c, identity)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	res := encodeLogin(user.Username, user.Name, user.Role, user.Active, accessToken, refreshToken)
	response.Default(c, res, "", http.StatusOK)
}

func (s service) HTTPRegister(c *gin.Context) {
	dr, err := decodeRegister(c)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	err = s.Register(c,
		model.Identity{
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
