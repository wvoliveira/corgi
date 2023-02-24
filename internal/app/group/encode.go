package group

import (
	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"time"
)

type addResponse struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

type listResponse struct {
	Groups []model.Group `json:"groups"`
	Limit  int           `json:"limit"`
	Page   int           `json:"page"`
	Sort   string        `json:"sort"`
	Total  int64         `json:"total"`
	Pages  int           `json:"pages"`
}

type findByIDUserModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type findByIDResponse struct {
	Group model.Group `json:"group"`
	Users []findByIDUserModel
}

type deleteResponse struct{}

type inviteModel struct {
	InviteID         string    `json:"invite_id"`
	CreatedAt        time.Time `json:"created_at"`
	GroupName        string    `json:"group_name"`
	GroupDisplayName string    `json:"group_display_name"`
	GroupDescription string    `json:"group_description"`
	InvitedByID      string    `json:"invited_by_id"`
	InvitedByName    string    `json:"invited_by_name"`
}

type invitesListResponse struct {
	Invites []inviteModel `json:"invites"`
	Page    int           `json:"page"`
	Pages   int           `json:"pages"`
	Total   int64         `json:"total"`
	Limit   int           `json:"limit"`
	Sort    string        `json:"sort"`
}

func encodeFindByID(c *gin.Context, group model.Group, users []model.User) (res findByIDResponse) {
	res.Group = group

	for _, user := range users {

		u := struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}{user.ID, user.Name}

		res.Users = append(res.Users, u)
	}

	return res
}

func encodeDelete(c *gin.Context) (res deleteResponse) {
	return
}
