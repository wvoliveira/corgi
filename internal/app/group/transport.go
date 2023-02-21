package group

import (
	"net/http"

	"github.com/gin-gonic/gin"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(rg *gin.RouterGroup) {
	r := rg.Group("/groups")

	r.POST("", s.HTTPAdd)
	r.GET("", s.HTTPList)
	r.GET("/:id", s.HTTPFindByID)
	r.DELETE("/:id", s.HTTPDelete)
	r.GET("/invites", s.HTTPInvitesList)
	r.POST("/:id/invites", s.HTTPInvitesAddByID)
	r.GET("/:id/invites", s.HTTPInvitesListByID)
}

func (s service) HTTPAdd(c *gin.Context) {
	d, err := decodeAdd(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	group := model.Group{Name: d.Name, DisplayName: d.DisplayName, Description: d.Description}
	group, err = s.Add(c, d.WhoID, group)
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
	d, err := decodeList(c)

	if err != nil {
		e.EncodeError(c, err)
		return
	}

	total, pages, groups, err := s.List(c, d.WhoID, d.Offset, d.Limit, d.Sort)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	resp := listResponse{
		Groups: groups,
		Limit:  d.Limit,
		Page:   d.Page,
		Sort:   d.Sort,
		Total:  total,
		Pages:  pages,
	}

	response.Default(c, resp, "", http.StatusOK)
}

func (s service) HTTPFindByID(c *gin.Context) {
	d, err := decodeFindByID(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	group, users, err := s.FindByID(c, d.WhoID, d.GroupID)
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

	err = s.Delete(c, d.WhoID, d.GroupID)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	// Not implemented yet.
	_ = encodeDelete(c)

	response.Default(c, nil, "", http.StatusOK)
}

func (s service) HTTPInvitesAddByID(c *gin.Context) {
	d, err := decodeInvitesAddByID(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	_, err = s.InvitesAddByID(c, d)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	response.Default(c, nil, "", http.StatusOK)
}

func (s service) HTTPInvitesListByID(c *gin.Context) {
	d, err := decodeInvitesListByID(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	total, pages, invites, err := s.InvitesListByID(c, d)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	resp := invitesListResponse{
		Invites: invites,
		Limit:   d.Limit,
		Page:    d.Page,
		Sort:    d.Sort,
		Total:   total,
		Pages:   pages,
	}

	response.Default(c, resp, "", http.StatusOK)
}

func (s service) HTTPInvitesList(c *gin.Context) {
	d, err := decodeInvitesList(c)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	total, pages, invites, err := s.InvitesList(c, d)
	if err != nil {
		e.EncodeError(c, err)
		return
	}

	resp := invitesListResponse{
		Invites: invites,
		Limit:   d.Limit,
		Page:    d.Page,
		Sort:    d.Sort,
		Total:   total,
		Pages:   pages,
	}

	response.Default(c, resp, "", http.StatusOK)
}
