package link

import (
	"context"
	"fmt"
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

func (s service) NatsBroker() {
	r.POST("", s.HTTPAdd)
	r.GET(":id", s.HTTPFindByID)
	r.GET("", s.HTTPFindAll)
	r.PATCH(":id", s.HTTPUpdate)
	r.DELETE(":id", s.HTTPDelete)
}

func (s service) NatsAdd() (*nats.Subscription, error) {
	sub, err := s.broker.Subscribe("link_add", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
		err := m.Respond([]byte("teste"))
		if err != nil {
			s.logger.With(context.TODO(), "error to respond to broker")
		}
	})
	if err != nil {
		return sub, err
	}
	return sub, nil
}

func (s service) NatsFindByID(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeFindByID(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	link, err := s.FindByID(c.Request.Context(), dr.ID, dr.UserID)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := findByIDResponse{
		Link: link,
		Err:  err,
	}
	encodeResponse(c, sr)
}

func (s service) NatsFindAll(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeFindAll(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	total, pages, links, err := s.FindAll(c.Request.Context(), dr.Offset, dr.Limit, dr.Sort, dr.UserID)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := findAllResponse{
		Links: links,
		Limit: dr.Limit,
		Page:  dr.Page,
		Sort:  dr.Sort,
		Total: total,
		Pages: pages,
		Err:   err,
	}
	encodeResponse(c, sr)
}

func (s service) NatsUpdate(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeUpdate(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	link, err := s.Update(
		c.Request.Context(),
		entity.Link{ID: dr.ID, Domain: dr.Domain, Keyword: dr.Keyword, URL: dr.URL, Title: dr.Title, Active: dr.Active, UserID: dr.UserID},
	)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := updateResponse{
		Link: link,
		Err:  err,
	}
	encodeResponse(c, sr)
}

func (s service) NatsDelete(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeDelete(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	err = s.Delete(c.Request.Context(), dr.ID, dr.UserID)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	sr := deleteResponse{
		Err: err,
	}
	encodeResponse(c, sr)
}
