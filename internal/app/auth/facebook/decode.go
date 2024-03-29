package facebook

import (
	"github.com/gin-gonic/gin"
)

type callbackRequest struct {
	State    string
	Code     string
	Scopes   []string
	AuthUser string
	Domain   string
	Prompt   string
}

func decodeCallbackRequest(c *gin.Context) (req callbackRequest, err error) {
	q := c.Request.URL.Query()
	var scopes []string
	req = callbackRequest{
		q.Get("state"),
		q.Get("code"),
		append(scopes, q.Get("scopes")),
		q.Get("authuser"),
		q.Get("hd"),
		q.Get("prompt"),
	}
	return
}
