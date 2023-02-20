package token

import (
	"github.com/gin-gonic/gin"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
)

func (s service) NewHTTP(rg *gin.RouterGroup) {
	r := rg.Group("/auth/token")

	r.GET("/valid", s.HTTPValid)
	r.POST("/refresh", s.HTTPRefresh)
}

func (s service) HTTPValid(c *gin.Context) {
	log := logger.Logger(c)

	err := e.ErrNotImplemented
	log.Warn().Caller().Msg(err.Error())
}

func (s service) HTTPRefresh(c *gin.Context) {
	log := logger.Logger(c)

	err := e.ErrNotImplemented
	log.Warn().Caller().Msg(err.Error())
}
