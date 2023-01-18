package log

import (
	"net/http"

	"github.com/gin-gonic/gin"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(rg *gin.RouterGroup) {
	r := rg.Group("/log")
	r.POST("", s.HTTPAdd)
}

func (s service) HTTPAdd(c *gin.Context) {
	user, payload, err := decodeAdd(c)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	err = s.Add(c, user, payload)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	resp := updateResponse{
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Err:       err,
	}

	response.Default(c, resp, "", http.StatusOK)
}
