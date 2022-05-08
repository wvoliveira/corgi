package redirect

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type findByKeywordRequest struct {
	LinkID  string `json:"link_id"`
	Domain  string `json:"domain"`
	Keyword string `json:"keyword"`
}

func decodeFindByKeyword(r *http.Request) (req findByKeywordRequest, err error) {
	domain := r.Host
	vars := mux.Vars(r)

	keyword := vars["keyword"]
	if keyword == "" {
		return req, errors.New("impossible to get redirect keyword")
	}

	req.Domain = domain
	req.Keyword = keyword
	return req, nil
}
