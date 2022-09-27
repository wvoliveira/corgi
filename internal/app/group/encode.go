package group

import "github.com/wvoliveira/corgi/internal/pkg/entity"

type addResponse struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Description string   `json:"description"`
	UserIDs     []string `json:"user_ids"`
}

type listResponse struct {
	Groups []entity.Group `json:"data"`
	Limit  int            `json:"limit"`
	Page   int            `json:"page"`
	Sort   string         `json:"sort"`
	Total  int64          `json:"total"`
	Pages  int            `json:"pages"`
}
