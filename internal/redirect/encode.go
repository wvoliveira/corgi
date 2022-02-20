package redirect

import (
	"github.com/elga-io/corgi/internal/entity"
)

type findByKeywordResponse struct {
	Link entity.Link `json:"data,omitempty"`
	Err  error       `json:"error,omitempty"`
}

func (r findByKeywordResponse) Error() error { return r.Err }
