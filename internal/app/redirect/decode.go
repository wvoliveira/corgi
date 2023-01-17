package redirect

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type findByKeywordRequest struct {
	LinkID  string `json:"link_id"`
	Domain  string `json:"domain"`
	Keyword string `json:"keyword"`
}

func decodeFindByKeyword(c *gin.Context) (r findByKeywordRequest, err error) {
	keyword := c.Param("keyword")

	if keyword == "" {
		return r, errors.New("impossible to get redirect keyword from path")
	}

	r.Keyword = keyword
	r.Domain = c.Request.Host

	return r, nil
}
