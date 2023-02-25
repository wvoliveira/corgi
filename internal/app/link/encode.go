package link

import "github.com/wvoliveira/corgi/internal/pkg/model"

type findRedirectResponse struct {
	URL string `json:"url"`
}

type findAllResponse struct {
	Links []model.Link `json:"links"`
	Limit int          `json:"limit"`
	Page  int          `json:"page"`
	Sort  string       `json:"sort"`
	Total int64        `json:"total"`
	Pages int          `json:"pages"`
	Err   error        `json:"error,omitempty"`
}

type findByKeywordResponse struct {
	URL string `json:"url"`
}

func encodeRedirect(link model.Link) (r findRedirectResponse) {
	r.URL = link.URL
	return
}

func encodeFindByKeyword(link model.Link) (r findByKeywordResponse) {
	r.URL = link.URL
	return
}
