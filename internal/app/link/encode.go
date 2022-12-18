package link

import "github.com/wvoliveira/corgi/internal/pkg/model"

type findAllResponse struct {
	Links []model.Link `json:"data"`
	Limit int          `json:"limit"`
	Page  int          `json:"page"`
	Sort  string       `json:"sort"`
	Total int64        `json:"total"`
	Pages int          `json:"pages"`
	Err   error        `json:"error,omitempty"`
}

func (r findAllResponse) Error() error { return r.Err }
