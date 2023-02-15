package token

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

func decodeLogout(c *gin.Context) (user model.User, err error) {
	log := logger.Logger(c.Request.Context())

	session := sessions.Default(c)
	v := session.Get("user")

	if v == nil {
		log.Warn().Caller().Msg("impossible to get user from session")
		return user, e.ErrUnauthorized
	}

	user = v.(model.User)

	// TODO: change to casbin or another authorization library.
	if user.ID == "anonymous" {
		log.Warn().Caller().Msg("impossible to get user from session")
		return user, e.ErrUnauthorized
	}
	return
}
