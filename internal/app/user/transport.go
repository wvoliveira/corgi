package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/middleware"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(rg *gin.RouterGroup) {
	r := rg.Group("/user")
	r.Use(middleware.Auth())

	r.GET("/me", s.HTTPFind)
	r.PATCH("/me", s.HTTPUpdate)
}

func (s service) HTTPFind(c *gin.Context) {

	var identities = []identity{}

	user, err := decodeFind(c)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	user, err = s.Find(c, user.ID)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	for _, i := range user.Identities {
		idt := identity{
			Provider: i.Provider,
			UID:      i.UID,
		}
		identities = append(identities, idt)
	}

	resp := userResponse{
		Name:       user.Name,
		Role:       user.Role,
		Identities: identities,
	}

	response.Default(c, resp, "", http.StatusOK)
}

func (s service) HTTPUpdate(c *gin.Context) {
	user, err := decodeUpdate(c)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	err = s.Update(c, model.User{ID: user.ID, Name: user.Name})

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	response.Default(c, nil, "", http.StatusOK)
}
