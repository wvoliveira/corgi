package user

import (
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func (s service) Routers(e *gin.Engine) {
	r := e.Group("/api/v1/user",
		middlewares.Access(s.logger),
		middlewares.Checks(s.logger),
		middlewares.Auth(s.logger, s.secret),
		middlewares.Authorizer(s.enforce, s.logger))

	r.GET("/me", s.HTTPFind)
	r.PATCH("/me", s.HTTPUpdate)
}

func (s service) HTTPFind(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeFind(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	user, err := s.Find(c.Request.Context(), dr.UserID)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	var identities []identity
	for _, i := range user.Identities {
		idt := identity{
			Provider: i.Provider,
			UID:      i.UID,
		}
		identities = append(identities, idt)
	}

	ur := userResponse{Name: user.Name, Role: user.Role, Identities: identities}
	sr := findResponse{
		userResponse: ur,
		Err:          err,
	}
	encodeResponse(c, sr)
}

func (s service) HTTPUpdate(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeUpdate(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	user, err := s.Update(c.Request.Context(), entity.User{ID: dr.UserID, Name: dr.Name})
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := updateResponse{
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Err:       err,
	}
	encodeResponse(c, sr)
}
