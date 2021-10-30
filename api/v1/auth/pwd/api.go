// Package pwd provides JSON Web Token (JWT) authentication and authorization middleware.
// It implements a passwordless authentication flow by sending login tokens vie email which are then exchanged for JWT access and refresh tokens.
package pwd

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"encoding/json"

	"github.com/go-chi/render"
	"github.com/go-kit/log"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/mssola/user_agent"

	"github.com/elga-io/redir/api/v1/auth/jwt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// AuthStorer defines database operations on accounts and tokens.
type AuthStorer interface {
	GetAccount(id int) (*Account, error)
	GetAccountByEmail(email string) (*Account, error)
	AuthAccount(email, password string) (*Account, error)
	UpdateAccount(a *Account) error

	GetToken(token string) (*jwt.Token, error)
	CreateOrUpdateToken(t *jwt.Token) error
	DeleteToken(t *jwt.Token) error
	PurgeExpiredToken() error
}

// Resource implements passwordless account authentication against a database.
type Resource struct {
	LoginAuth *LoginTokenAuth
	TokenAuth *jwt.TokenAuth
	Store     AuthStorer
	logger    log.Logger
}

// NewResource returns a configured authentication resource.
func NewResource(authStore AuthStorer) (*Resource, error) {
	loginAuth, err := NewLoginTokenAuth()
	if err != nil {
		return nil, err
	}

	tokenAuth, err := jwt.NewTokenAuth()
	if err != nil {
		return nil, err
	}

	resource := &Resource{
		LoginAuth: loginAuth,
		TokenAuth: tokenAuth,
		Store:     authStore,
	}

	resource.choresTicker()

	return resource, nil
}

// Router provides necessary routes for password authentication flow.
func (rs *Resource) Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/login", rs.login).Methods("POST")
	r.HandleFunc("/token", rs.token).Methods("POST")

	s := r.NewRoute().Subrouter()
	s.Use(rs.TokenAuth.Verifier())
	s.Use(jwt.AuthenticateRefreshJWT)
	s.HandleFunc("/refresh", rs.refresh).Methods("POST")
	s.HandleFunc("/logout", rs.logout).Methods("POST")
	return r
}

type loginRequest struct {
	Email    string `json:"email" gorm:"email"`
	Password string `json:"password" gorm:"password"`
}

func (body *loginRequest) Bind(r *http.Request) error {
	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	return validation.ValidateStruct(body,
		validation.Field(&body.Email, validation.Required, is.Email),
		validation.Field(&body.Password, validation.Required, is.ASCII),
	)
}

func (rs *Resource) login(w http.ResponseWriter, r *http.Request) {
	lr := &loginRequest{}

	err := json.NewDecoder(r.Body).Decode(&lr)
	if err != nil {
		rs.logger.Log("email", lr.Email, "warn", err)
		fmt.Fprintln(w, ErrUnauthorized(ErrInvalidLogin))
		return
	}

	acc, err := rs.Store.GetAccountByEmail(lr.Email)
	if err != nil {
		rs.logger.Log("email", lr.Email, "warn", err)
		fmt.Fprintln(w, ErrUnauthorized(ErrUnknownLogin))
		return
	}

	if !*acc.CanLogin() {
		fmt.Fprintln(w, ErrUnauthorized(ErrLoginDisabled))
		return
	}

	ua := user_agent.New(r.UserAgent())
	browser, _ := ua.Browser()

	token := &jwt.Token{
		Token:      uuid.Must(uuid.NewV4()).String(),
		Expiry:     time.Now().Add(rs.TokenAuth.JwtRefreshExpiry),
		UpdatedAt:  time.Now(),
		AccountID:  acc.ID,
		Mobile:     ua.Mobile(),
		Identifier: fmt.Sprintf("%s on %s", browser, ua.OS()),
	}

	if err := rs.Store.CreateOrUpdateToken(token); err != nil {
		rs.logger.Log("method", "Store.CreateOrUpdateToken", "err", err)
		fmt.Fprintln(w, ErrInternalServerError)
		return
	}

	access, refresh, err := rs.TokenAuth.GenTokenPair(acc.Claims(), token.Claims())
	if err != nil {
		rs.logger.Log("method", "TokenAuth.GenTokenPair", "err", err)
		fmt.Fprintln(w, ErrInternalServerError)
		return
	}

	acc.LastLogin = time.Now()
	if err := rs.Store.UpdateAccount(acc); err != nil {
		rs.logger.Log("method", "Store.UpdateAccount", "err", err)
		fmt.Fprintln(w, ErrInternalServerError)
		return
	}

	render.Respond(w, r, &tokenResponse{
		Access:        access,
		Refresh:       refresh,
		AccessExpiry:  int(rs.TokenAuth.JwtExpiry.Seconds()),
		RefreshExpiry: int(rs.TokenAuth.JwtRefreshExpiry.Seconds()),
	})
}

type tokenRequest struct {
	Token string `json:"token"`
}

type tokenResponse struct {
	Access        string `json:"access_token"`
	Refresh       string `json:"refresh_token"`
	AccessExpiry  int    `json:"expires_in"`
	RefreshExpiry int    `json:"refresh_token_expires_in"`
}

func (body *tokenRequest) Bind(r *http.Request) error {
	body.Token = strings.TrimSpace(body.Token)

	return validation.ValidateStruct(body,
		validation.Field(&body.Token, validation.Required, is.Alphanumeric),
	)
}

func (rs *Resource) token(w http.ResponseWriter, r *http.Request) {
	tr := &tokenRequest{}
	err := json.NewDecoder(r.Body).Decode(&tr)
	if err != nil {
		rs.logger.Log("method", "NewDecoder", "warn", err)
		fmt.Fprintln(w, ErrUnauthorized(ErrLoginToken))
		return
	}

	id, err := rs.LoginAuth.GetAccountID(tr.Token)
	if err != nil {
		fmt.Fprintln(w, ErrUnauthorized(ErrLoginToken))
		return
	}

	acc, err := rs.Store.GetAccount(id)
	if err != nil {
		// account deleted before login token expired
		fmt.Fprintln(w, ErrUnauthorized(ErrUnknownLogin))
		return
	}

	if !*acc.CanLogin() {
		fmt.Fprintln(w, ErrUnauthorized(ErrLoginDisabled))
		return
	}

	ua := user_agent.New(r.UserAgent())
	browser, _ := ua.Browser()

	token := &jwt.Token{
		Token:      uuid.Must(uuid.NewV4()).String(),
		Expiry:     time.Now().Add(rs.TokenAuth.JwtRefreshExpiry),
		UpdatedAt:  time.Now(),
		AccountID:  acc.ID,
		Mobile:     ua.Mobile(),
		Identifier: fmt.Sprintf("%s on %s", browser, ua.OS()),
	}

	if err := rs.Store.CreateOrUpdateToken(token); err != nil {
		rs.logger.Log("method", "Store.CreateOrUpdateToken", "err", err)
		fmt.Fprintln(w, ErrInternalServerError)
		return
	}

	access, refresh, err := rs.TokenAuth.GenTokenPair(acc.Claims(), token.Claims())
	if err != nil {
		rs.logger.Log("method", "TokenAuth.GenTokenPair", "err", err)
		fmt.Fprintln(w, ErrInternalServerError)
		return
	}

	acc.LastLogin = time.Now()
	if err := rs.Store.UpdateAccount(acc); err != nil {
		rs.logger.Log("method", "Store.UpdateAccount", "err", err)
		fmt.Fprintln(w, ErrInternalServerError)
		return
	}

	render.Respond(w, r, &tokenResponse{
		Access:  access,
		Refresh: refresh,
	})
}

func (rs *Resource) refresh(w http.ResponseWriter, r *http.Request) {
	rt := jwt.RefreshTokenFromCtx(r.Context())

	token, err := rs.Store.GetToken(rt)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(jwt.ErrTokenExpired))
		return
	}

	if time.Now().After(token.Expiry) {
		rs.Store.DeleteToken(token)
		render.Render(w, r, ErrUnauthorized(jwt.ErrTokenExpired))
		return
	}

	acc, err := rs.Store.GetAccount(token.AccountID)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(ErrUnknownLogin))
		return
	}

	if !*acc.CanLogin() {
		render.Render(w, r, ErrUnauthorized(ErrLoginDisabled))
		return
	}

	token.Token = uuid.Must(uuid.NewV4()).String()
	token.Expiry = time.Now().Add(rs.TokenAuth.JwtRefreshExpiry)
	token.UpdatedAt = time.Now()

	access, refresh, err := rs.TokenAuth.GenTokenPair(acc.Claims(), token.Claims())
	if err != nil {
		rs.logger.Log("method", "TokenAuth.GenTokenPair", "err", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	if err := rs.Store.CreateOrUpdateToken(token); err != nil {
		rs.logger.Log("method", "Store.CreateOrUpdateToken", "err", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	acc.LastLogin = time.Now()
	if err := rs.Store.UpdateAccount(acc); err != nil {
		rs.logger.Log("method", "Store.UpdateAccount", "err", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Respond(w, r, &tokenResponse{
		Access:  access,
		Refresh: refresh,
	})
}

func (rs *Resource) logout(w http.ResponseWriter, r *http.Request) {
	rt := jwt.RefreshTokenFromCtx(r.Context())
	token, err := rs.Store.GetToken(rt)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(jwt.ErrTokenExpired))
		return
	}
	rs.Store.DeleteToken(token)

	render.Respond(w, r, http.NoBody)
}
