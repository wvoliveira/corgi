package redirect

import (
	"net/http"

	e "github.com/elga-io/corgi/internal/pkg/errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (s service) NewHTTP(e *gin.Engine) {
	r := e.Group("/")

	r.Use(sessions.Sessions("_corgi", s.store))
	r.GET("/:keyword", s.HTTPFind)
}

func (s service) HTTPFind(c *gin.Context) {
	// Decode request to request object.
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

	// Pass decode request to from gin context to use in middleware.
	c.Set("findByKeywordResponse", link)

	// Redirect! Not encode for response.
	c.Redirect(http.StatusMovedPermanently, link.URL)
}
