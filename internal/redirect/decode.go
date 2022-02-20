package redirect

import (
	"github.com/gin-gonic/gin"
)

type findByKeywordRequest struct {
	LinkID  string `json:"link_id"`
	Domain  string `json:"domain"`
	Keyword string `json:"keyword"`
}

func decodeFindByKeyword(c *gin.Context) (req findByKeywordRequest, err error) {
	domain := c.Request.Host
	keyword := c.Param("keyword")

	req.Domain = domain
	req.Keyword = keyword
	return req, nil
}
