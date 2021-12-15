package public

import (
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/middlewares"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ua-parser/uap-go/uaparser"
	"net/http"
)

func (s service) Routers(e *gin.Engine) {
	r := e.Group("/",
		middlewares.Access(s.logger),
		sessions.SessionsMany([]string{"_corgi", "session"}, s.store))

	r.GET("/:keyword", s.HTTPFindByKeyword)
}

func (s service) HTTPFindByKeyword(c *gin.Context) {
	remoteAddress := c.ClientIP()
	userAgent := c.Request.Header.Get("user-agent")
	referer := c.Request.Header.Get("referer")

	parser := uaparser.NewFromSaved()
	client := parser.Parse(userAgent)
	// n := net.ParseIP(remoteAddress)

	// Decode request to request object.
	dr, err := decodeFindByKeyword(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Check if key for this keyword exists.
	unique := true
	var keywords []string

	sessionUnique := sessions.DefaultMany(c, "_corgi")
	domainKey := sessionUnique.Get(dr.Domain)

	if domainKey != nil {
		keywords = domainKey.([]string)
		for _, keyword := range keywords {
			if keyword == dr.Keyword {
				unique = false
				break
			}
		}
	}

	linkLog := entity.LinkLog{
		RemoteAddress:         remoteAddress, // change to remoteAddress
		UserAgentRaw:          userAgent,
		UserAgentFamily:       client.UserAgent.Family,
		UserAgentOSFamily:     client.Os.Family,
		UserAgentDeviceFamily: client.Device.Family,
		Referer:               referer,
	}

	// Business logic.
	link, err := s.FindByKeyword(c.Request.Context(), dr.Domain, dr.Keyword, linkLog, unique)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Set keyword in session domain key
	if unique {
		keywords = append(keywords, dr.Keyword)
		sessionUnique.Set(dr.Domain, keywords)
		_ = sessionUnique.Save()
	}

	// Redirect! Not encode for response.
	sr := findByKeywordResponse{
		Link: link,
		Err:  err,
	}
	c.Redirect(http.StatusMovedPermanently, sr.Link.URL)
}
