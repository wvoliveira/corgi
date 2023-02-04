package redirect

import (
	"net/http"

	"github.com/gin-gonic/gin"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(rg *gin.RouterGroup) {
	r := rg.Group("/redirect")
	// r.Use(middleware.UniqueUserForKeywords())

	r.GET("/:keyword", s.HTTPFind)
}

func (s service) HTTPFind(c *gin.Context) {

	dr, err := decodeFindByKeyword(c)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	link, err := s.Find(c, dr.Domain, dr.Keyword)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	response.Default(c, link, "", http.StatusOK)
}
