package link

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wvoliveira/corgi/internal/app/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/middleware"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(r *mux.Router) {
	rr := r.PathPrefix("/v1/links").Subrouter()
	rr.Use(middleware.Checks)
	rr.Use(middleware.Auth(s.secret))

	rr.HandleFunc("", nil).Methods("OPTIONS")
	rr.HandleFunc("", s.HTTPAdd).Methods("POST")
	rr.HandleFunc("/:id", s.HTTPFindByID).Methods("GET")
	rr.HandleFunc("/status/:id", s.HTTPFindByID).Methods("GET")
	rr.HandleFunc("", s.HTTPFindAll).Methods("GET")
	rr.HandleFunc("/:id", s.HTTPUpdate).Methods("PATCH")
	rr.HandleFunc("/:id", s.HTTPDelete).Methods("DELETE")
}

func (s service) HTTPAdd(w http.ResponseWriter, r *http.Request) {
	dr, err := decodeAdd(r)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	link, err := s.Add(r.Context(), entity.Link{Domain: dr.Domain, Keyword: dr.Keyword, URL: dr.URL, Title: dr.Title, UserID: dr.UserID})
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	response.Default(w, link, "", http.StatusCreated)
}

func (s service) HTTPFindByID(w http.ResponseWriter, r *http.Request) {
	dr, err := decodeFindByID(r)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	link, err := s.FindByID(r.Context(), dr.ID, dr.UserID)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	response.Default(w, link, "", http.StatusOK)
}

func (s service) HTTPFindAll(w http.ResponseWriter, r *http.Request) {
	dr, err := decodeFindAll(r)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	total, pages, links, err := s.FindAll(r.Context(), dr.Offset, dr.Limit, dr.Sort, dr.UserID)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	sr := findAllResponse{
		Links: links,
		Limit: dr.Limit,
		Page:  dr.Page,
		Sort:  dr.Sort,
		Total: total,
		Pages: pages,
		Err:   err,
	}

	response.Default(w, sr, "", http.StatusOK)
}

func (s service) HTTPUpdate(w http.ResponseWriter, r *http.Request) {
	dr, err := decodeUpdate(r)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	link, err := s.Update(r.Context(), entity.Link{
		ID:      dr.ID,
		Domain:  dr.Domain,
		Keyword: dr.Keyword,
		URL:     dr.URL,
		Title:   dr.Title,
		Active:  dr.Active,
		UserID:  dr.UserID,
	})
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	response.Default(w, link, "", http.StatusOK)
}

func (s service) HTTPDelete(w http.ResponseWriter, r *http.Request) {
	dr, err := decodeDelete(r)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	err = s.Delete(r.Context(), dr.ID, dr.UserID)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	response.Default(w, nil, "", http.StatusOK)
}
