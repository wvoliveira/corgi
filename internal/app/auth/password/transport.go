package password

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
)

func (s service) NewHTTP(rg *gin.RouterGroup) {
	r := rg.Group("/auth/password")
	// r.Use(middleware.Checks)

	r.POST("/login", s.HTTPLogin())
	r.POST("/register", s.HTTPRegister())
}

func (s service) HTTPLogin() gin.HandlerFunc {
	return func(c *gin.Context) {

		dr, err := decodeLogin(c)
		if err != nil {
			e.EncodeError(c, err)
			return
		}

		identity := entity.Identity{
			Provider: "email",
			UID:      dr.Email,
			Password: dr.Password,
		}

		tokenAccess, tokenRefresh, err := s.Login(c, identity)
		if err != nil {
			e.EncodeError(c, err)
			return
		}

		// TODOS:
		// - change token model to create access and refresh token only
		// - use options to set Domain and MaxAge. Ex.:
		//	 https://github.com/gin-contrib/sessions/blob/master/session_options_go1.10.go
		session := sessions.Default(c)
		session.Set("token_access", tokenAccess.Token)
		session.Set("token_refresh", tokenRefresh.ID)

		c.String(http.StatusOK, "")
	}
}

func (s service) HTTPRegister() gin.HandlerFunc {
	return func(c *gin.Context) {

		dr, err := decodeRegister(c)
		if err != nil {
			e.EncodeError(c, err)
			return
		}

		err = s.Register(c, entity.Identity{
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
}
