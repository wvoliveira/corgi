package server

import "net/http"

func getAccountFromHeaders(r *http.Request) (a Account) {
	a.ID = r.Header.Get("AccountID")
	a.Email = r.Header.Get("AccountEmail")
	a.Role = r.Header.Get("AccountRole")
	return
}
