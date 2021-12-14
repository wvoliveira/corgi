package public

import (
	"github.com/gin-gonic/gin"
)

type findByKeywordRequest struct {
	Domain  string `json:"-"`
	Keyword string `json:"-"`
}

func decodeFindByKeyword(c *gin.Context) (req findByKeywordRequest, err error) {
	domain := c.Request.Host
	keyword := c.Param("keyword")

	req.Domain = domain
	req.Keyword = keyword
	return req, nil
}
