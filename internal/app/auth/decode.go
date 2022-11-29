package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
)

func decodeLogout(c *gin.Context) (user entity.User, err error) {
	log := logger.Logger(c.Request.Context())

	session := sessions.Default(c)
	v := session.Get("user")

	if v == nil {
		log.Warn().Caller().Msg("impossible to get user from session")
		return user, e.ErrUnauthorized
	}

	user = v.(entity.User)

	// TODO: change to casbin or another authorization library.
	if user.ID == "anonymous" {
		log.Warn().Caller().Msg("impossible to get user from session")
		return user, e.ErrUnauthorized
	}
	return
}
