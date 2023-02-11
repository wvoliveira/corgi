package short

import (
	"net/http"

	"github.com/gin-gonic/gin"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(rg *gin.RouterGroup) {
	r := rg.Group("/short")
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

	url := encodeFindByKeyword(link)

	response.Default(c, url, "", http.StatusOK)
}
