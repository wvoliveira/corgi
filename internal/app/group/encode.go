package group

import "github.com/wvoliveira/corgi/internal/pkg/model"

type addResponse struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

type listResponse struct {
	Groups []model.Group `json:"data"`
	Limit  int           `json:"limit"`
	Page   int           `json:"page"`
	Sort   string        `json:"sort"`
	Total  int64         `json:"total"`
	Pages  int           `json:"pages"`
}
