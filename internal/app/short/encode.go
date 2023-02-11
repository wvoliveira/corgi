package short

import (
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

type findByKeywordResponse struct {
	URL string `json:"url"`
}

func encodeFindByKeyword(link model.Link) (r findByKeywordResponse) {
	r.URL = link.URL
	return
}
