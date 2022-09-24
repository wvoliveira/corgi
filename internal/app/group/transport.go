package group

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/middleware"
	"github.com/wvoliveira/corgi/internal/pkg/response"
)

func (s service) NewHTTP(r *mux.Router) {
	rr := r.PathPrefix("/v1/groups").Subrouter()
	rr.Use(middleware.Auth(s.secret))

	rr.HandleFunc("", s.HTTPAdd).Methods("POST")
	rr.HandleFunc("", s.HTTPList).Methods("GET")
}

func (s service) HTTPAdd(w http.ResponseWriter, r *http.Request) {
	payload, userID, err := decodeAdd(r)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	group := entity.Group{Name: payload.Name, Description: payload.Description}

	group, err = s.Add(r.Context(), group, userID)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	response.Default(w, group, "", http.StatusOK)
}

func (s service) HTTPList(w http.ResponseWriter, r *http.Request) {
	payload, err := decodeList(r)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	total, pages, groups, err := s.List(r.Context(), payload.Offset, payload.Limit, payload.Sort, payload.UserID)
	if err != nil {
		e.EncodeError(w, err)
		return
	}

	resp := listResponse{
		Groups: groups,
		Limit:  payload.Limit,
		Page:   payload.Page,
		Sort:   payload.Sort,
		Total:  total,
		Pages:  pages,
	}

	response.Default(w, resp, "", http.StatusOK)
}
