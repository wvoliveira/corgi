package group

import (
	"net/http"

	"github.com/gin-gonic/gin"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/middleware"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(rg *gin.RouterGroup) {
	r := rg.Group("/groups")
	r.Use(middleware.Auth())

	r.POST("", s.HTTPAdd)
	r.GET("", s.HTTPList)
	r.GET("/:id", s.HTTPFindByID)
	r.DELETE("/:id", s.HTTPDelete)
}

func (s service) HTTPAdd(c *gin.Context) {
	payload, userID, err := decodeAdd(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	group := model.Group{Name: payload.Name, DisplayName: payload.DisplayName, Description: payload.Description}

	group, err = s.Add(c, group, userID)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	resp := addResponse{
		Name:        group.Name,
		DisplayName: group.DisplayName,
		Description: group.Description,
	}

	response.Default(c, resp, "", http.StatusOK)
}

func (s service) HTTPList(c *gin.Context) {
	payload, userID, err := decodeList(c)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	total, pages, groups, err := s.List(c, payload.Offset, payload.Limit, payload.Sort, userID)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	resp := listResponse{
		Groups: groups,
		Limit:  payload.Limit,
		Page:   payload.Page,
		Sort:   payload.Sort,
		Total:  total,
		Pages:  pages,
	}

	response.Default(c, resp, "", http.StatusOK)
}

func (s service) HTTPFindByID(c *gin.Context) {
	payload, userID, err := decodeFindByID(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	group, users, err := s.FindByID(c, payload.ID, userID)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	resp := encodeFindByID(c, group, users)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	response.Default(c, resp, "", http.StatusOK)
}

func (s service) HTTPDelete(c *gin.Context) {
	d, err := decodeDelete(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	err = s.Delete(c, d.UserID, d.GroupID)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Not implemented yet.
	_ = encodeDelete(c)

	response.Default(c, nil, "", http.StatusOK)
}
