package redirect

import (
	"context"
	"fmt"
	"github.com/elga-io/corgi/internal/entity"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/elga-io/corgi/pkg/queue"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ua-parser/uap-go/uaparser"
	"net"
	"time"
)

// MiddlewareMetric returns a middleware that records a metric in AMQP for each successful redirect.
func (s service) MiddlewareMetric(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		ctx := c.Request.Context()

		l := logger.With(ctx)

		remoteAddress := c.ClientIP()
		userAgent := c.Request.Header.Get("user-agent")
		referer := c.Request.Header.Get("referer")

		parser := uaparser.NewFromSaved()
		client := parser.Parse(userAgent)
		ip := net.ParseIP(remoteAddress)

		// Check if key for this keyword exists.
		unique := true
		var keywords []string

		dr := findByKeywordResponse{}

		v, exist := c.Get("findByKeywordResponse")
		if exist {
			dr = v.(findByKeywordResponse)
		} else {
			fmt.Println("nao existe")
			return
		}

		session := sessions.Default(c)
		domainKey := session.Get(dr.Link.Domain)

		fmt.Println("domain key:")
		fmt.Println(domainKey)

		// Check unique cookie session.
		// If it's nil or doesn't exist, this is a first click.
		if domainKey != nil {
			keywords = domainKey.([]string)
			for _, keyword := range keywords {
				fmt.Println("keyword:")
				fmt.Println(keyword)

				if keyword == dr.Link.Keyword {
					unique = false
					break
				}
			}
		}

		linkLog := entity.LinkLog{
			CreatedAt:             time.Now(),
			RemoteAddress:         ip.String(),
			UserAgent:             userAgent,
			UserAgentFamily:       client.UserAgent.Family,
			UserAgentOSFamily:     client.Os.Family,
			UserAgentDeviceFamily: client.Device.Family,
			Referer:               referer,
			LinkID:                dr.Link.ID,

			// Log object.
			Domain:  dr.Link.Domain,
			Keyword: dr.Link.Keyword,
			URL:     dr.Link.URL,
			Title:   dr.Link.Title,
		}

		// Set keyword in session domain key
		if !unique {
			fmt.Println("nao eh unico")
			return
		}

		keywords = append(keywords, dr.Link.Keyword)
		session.Set(dr.Link.Domain, keywords)
		err := session.Save()
		if err != nil {
			l.Warnf("error to save session with redirect domain: %s", err.Error())
		}

		fmt.Println("keywords")
		fmt.Println(keywords)

		fmt.Println("midd")
		fmt.Println(linkLog)

		go func() {
			// body, err := json.Marshal(linkLog)
			//if err != nil {
			//	fmt.Printf("error to marshal link log: %s\n\n", err.Error())
			//}

			data := queue.SendRequest{
				Body:     "Testing 1,2,3,...",                   // Required
				QueueURL: "http://localhost:9324/queue/default", // Required
			}

			resp, err := s.queuer.Send(context.TODO(), &data)

			if err != nil {
				l.Warnf("error to send message to MQ: %s", err.Error())
			}

			fmt.Printf("id returned by send request to MQ: %s\n\n", resp)
		}()

		c.Next()

		session = sessions.Default(c)
		domainKey = session.Get(dr.Link.Domain)

		fmt.Println("domain key:")
		fmt.Println(domainKey)

		err = session.Save()
		if err != nil {
			l.Warnf("error to save session with redirect domain: %s", err.Error())
		}

	}
}
