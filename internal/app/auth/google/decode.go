package google

import "github.com/gin-gonic/gin"

type callbackRequest struct {
	State    string   //state=state
	Code     string   // code=4/0AX4XfWjLX8K0dMCvLgUA86jMy_nTuRhuAKLtxGSC0gFvD8xWiNx-JjEDZ-XX4c93Wq1wzg
	Scopes   []string // scope=email%20profile%20https://www.googleapis.com/auth/userinfo.email%20https://www.googleapis.com/auth/userinfo.profile%20openid
	AuthUser string   // authuser=0
	Domain   string   // hd = elga.io
	Prompt   string   //prompt = consent
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
