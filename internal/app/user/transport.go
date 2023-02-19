package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/middleware"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(rg *gin.RouterGroup) {
	r := rg.Group("/users")
	r.Use(middleware.Auth())

	r.GET("/me", s.HTTPFindMe)
	r.PATCH("/me", s.HTTPUpdateMe)
	r.GET("/:id_username", s.HTTPFindByIDorUsername)
	r.PATCH("/:id_username", s.HTTPUpdateByIDorUsername)
}

func (s service) HTTPFindMe(c *gin.Context) {
	var identities = []identity{}

	d, err := decodeFindMe(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	user, err := s.FindMe(c, d.whoID)
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
		Username:   user.Username,
		Name:       user.Name,
		Role:       user.Role,
		Identities: identities,
	}

	response.Default(c, resp, "", http.StatusOK)
}

func (s service) HTTPUpdateMe(c *gin.Context) {
	d, err := decodeUpdateMe(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	err = s.UpdateMe(c, d.whoID, d.User.Name)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	response.Default(c, nil, "", http.StatusOK)
}

func (s service) HTTPFindByIDorUsername(c *gin.Context) {
	var identities = []identity{}

	d, err := decodeFindByIDorUsername(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	user, err := s.FindByIDorUsername(c, d.whoID, d.IDorUsername)
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
		Username:   user.Username,
		Name:       user.Name,
		Role:       user.Role,
		Identities: identities,
	}

	response.Default(c, resp, "", http.StatusOK)
}

func (s service) HTTPUpdateByIDorUsername(c *gin.Context) {
	d, err := decodeUpdateByIDorUsername(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	err = s.UpdateByIDorUsername(c, d.whoID, d.IDorUsername, d.User.Name)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	response.Default(c, nil, "", http.StatusOK)
}
