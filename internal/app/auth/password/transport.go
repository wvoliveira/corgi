package password

import (
	"encoding/json"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
)

func (s service) NewHTTP(rg *gin.RouterGroup) {
	r := rg.Group("/auth/password")
	// r.Use(middleware.Checks)

	r.POST("/login", s.HTTPLogin)
	r.POST("/register", s.HTTPRegister)
}

func (s service) HTTPLogin(c *gin.Context) {

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

	user, err := s.Login(c, identity)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	userAsText, err := json.Marshal(user)
	if err != nil {
		e.EncodeError(c, err)
	}

	// TODOS:
	// - change token model to create access and refresh token only
	// - use options to set Domain and MaxAge. Ex.:
	//	 https://github.com/gin-contrib/sessions/blob/master/session_options_go1.10.go
	session := sessions.Default(c)
	session.Set("user", string(userAsText))

	err = session.Save()
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	fmt.Println(session.Get("user"))

	c.JSON(200, gin.H{
		"name":   user.Name,
		"role":   user.Role,
		"active": user.Active,
	})
}

func (s service) HTTPRegister(c *gin.Context) {

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
