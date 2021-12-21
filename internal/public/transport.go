package public

import (
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ua-parser/uap-go/uaparser"
	"net"
	"net/http"
)

func (s service) Routers(e *gin.Engine) {
	r := e.Group("/",
		sessions.Sessions("_corgi", s.store))

	r.GET("/:keyword", s.HTTPFindByKeyword)
}

func (s service) HTTPFindByKeyword(c *gin.Context) {
	remoteAddress := c.ClientIP()
	userAgent := c.Request.Header.Get("user-agent")
	referer := c.Request.Header.Get("referer")

	parser := uaparser.NewFromSaved()
	client := parser.Parse(userAgent)
	ip := net.ParseIP(remoteAddress)

	// Decode request to request object.
	dr, err := decodeFindByKeyword(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Check if key for this keyword exists.
	unique := true
	var keywords []string

	session := sessions.Default(c)
	domainKey := session.Get(dr.Domain)

	if ip.IsLoopback() {
		unique = false
	}

	if unique && domainKey != nil {
		keywords = domainKey.([]string)
		for _, keyword := range keywords {
			if keyword == dr.Keyword {
				unique = false
				break
			}
		}
	}

	linkLog := entity.LinkLog{
		RemoteAddress:         remoteAddress,
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
		session.Set(dr.Domain, keywords)
		_ = session.Save()
	}

	// Redirect! Not encode for response.
	sr := findByKeywordResponse{
		Link: link,
		Err:  err,
	}
	c.Redirect(http.StatusMovedPermanently, sr.Link.URL)
}
