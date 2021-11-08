package server

import (
	"os"

	"github.com/gorilla/sessions"
)

var (
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
)
