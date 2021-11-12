package server

import (
	"os"

	"github.com/gorilla/sessions"
)

var (
	store     = sessions.NewCookieStore([]byte(os.Getenv("REDIR_SESSION_KEY")))
	secretKey = os.Getenv("REDIR_SECRET_KEY")
)
