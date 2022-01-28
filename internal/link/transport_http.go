package link

import (
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"time"
)

func (s service) HTTPNewTransport(e *gin.Engine) {
	r := e.Group("/api/v1/links",
		middlewares.Checks(s.logger),
		middlewares.Auth(s.logger, s.secret))

	r.OPTIONS("", nil)
	r.POST("", s.HTTPAdd)
	r.GET(":id", s.HTTPFindByID)
	r.GET("/status/:id", s.HTTPFindByID)
	r.GET("", s.HTTPFindAll)
	r.PATCH(":id", s.HTTPUpdate)
	r.DELETE(":id", s.HTTPDelete)
}

func (s service) HTTPAdd(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeAdd(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	request := dr
	response := addResponse{}

	// Get data from broker.
	// TODO: get error from NATS response.
	err = s.broker.Request("link.add", request, &response, 5*time.Second)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	encodeResponse(c, response)
}

func (s service) HTTPFindByID(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeFindByID(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	request := dr
	response := findByIDResponse{}

	// Get data from broker.
	err = s.broker.Request("link.findbyid", request, &response, 5*time.Second)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	encodeResponse(c, response)
}

func (s service) HTTPFindAll(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeFindAll(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	request := dr
	response := findAllResponse{}

	// Get data from broker.
	// TODO: get error from NATS response.
	err = s.broker.Request("link.findall", request, &response, 5*time.Second)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	encodeResponse(c, response)
}

func (s service) HTTPUpdate(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeUpdate(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	request := dr
	response := updateResponse{}

	// Get data from broker.
	err = s.broker.Request("link.update", request, &response, 5*time.Second)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	encodeResponse(c, response)
}

func (s service) HTTPDelete(c *gin.Context) {
	// Decode request to request object.
	dr, err := decodeDelete(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Business logic.
	request := dr
	response := deleteResponse{}

	// Get data from broker.
	err = s.broker.Request("link.delete", request, &response, 5*time.Second)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Encode object to answer request (response).
	encodeResponse(c, response)
}
