package redirect

import (
	"net/http"

	"github.com/gorilla/mux"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/middleware"
)

func (s service) NewHTTP(r *mux.Router) {
	rr := r.PathPrefix("/").Subrouter()
	rr.Use(middleware.SesssionRedirect(s.store, "_corgi"))

	rr.HandleFunc("/:keyword", s.HTTPFind).Methods("GET")
}

func (s service) HTTPFind(w http.ResponseWriter, r *http.Request) {
	// Decode request to request object.
	dr, err := decodeFindByKeyword(r)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	link, err := s.Find(r.Context(), dr.Domain, dr.Keyword)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	// Pass decode request to from gin context to use in middleware.
	// c.Set("findByKeywordResponse", link)

	// Redirect! Not encode for response.
	http.Redirect(w, r, link.URL, http.StatusMovedPermanently)
}
