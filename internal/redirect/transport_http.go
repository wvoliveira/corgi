package redirect

import (
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (s service) HTTPNewTransport(e *gin.Engine) {
	r := e.Group("/")

	r.Use(sessions.Sessions("_corgi", s.store))
	r.Use(s.MiddlewareMetric(s.logger))
	r.GET("/:keyword", s.HTTPFindByKeyword)
}

func (s service) HTTPFindByKeyword(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeFindByKeyword(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	request := dr
	response := findByKeywordResponse{}

	// Business logic.
	// Get data from broker.
	// TODO: get error from NATS response.
	err = s.broker.Request("link.findbykeyword", request, &response, 5*time.Second)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Pass decode request to from gin context to use in middleware.
	c.Set("findByKeywordResponse", response)

	// Redirect! Not encode for response.
	c.Redirect(http.StatusMovedPermanently, response.Link.URL)
}
