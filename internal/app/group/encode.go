package group

import (
	"github.com/gin-gonic/gin"
	"github.com/wvoliveira/corgi/internal/pkg/model"
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

type findByIDResponse struct {
	Group model.Group `json:"group"`
	// TODO: melhorar isso.. pelo amor
	// Preciso realizar o encode de somente alguns valores do usuário
	Users []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
}

type deleteResponse struct{}

type invitesListResponse struct {
	Invites []model.GroupInvite `json:"invites"`
	Limit   int                 `json:"limit"`
	Page    int                 `json:"page"`
	Sort    string              `json:"sort"`
	Total   int64               `json:"total"`
	Pages   int                 `json:"pages"`
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
