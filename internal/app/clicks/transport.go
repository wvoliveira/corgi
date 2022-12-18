package clicks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/middleware"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(rg *gin.RouterGroup) {
	r := rg.Group("/clicks")
	r.Use(middleware.Checks())
	r.Use(middleware.Auth())

	r.GET("/:link", s.HTTPFind)

}

func (s service) HTTPFind(c *gin.Context) {

	link := decodeFind(c)

	data, err := s.Find(c, link)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	response.Default(c, data, "", http.StatusOK)
}
