package group

import "github.com/wvoliveira/corgi/internal/pkg/entity"

type listResponse struct {
	Groups []entity.Group `json:"data"`
	Limit  int            `json:"limit"`
	Page   int            `json:"page"`
	Sort   string         `json:"sort"`
	Total  int64          `json:"total"`
	Pages  int            `json:"pages"`
}
