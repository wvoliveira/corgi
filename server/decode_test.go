package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestDecodeAuthLoginRequest(t *testing.T) {
	payloads := []authLoginRequest{
		{Email: "test1@email.com", Password: "secret1"},
		{Email: "teste2@email.com", Password: "secret2"},
	}

	for _, payload := range payloads {
		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewReader(body))
		if err != nil {
			t.Errorf("error to create a new request to /signin: %s", err.Error())
		}

		dr, err := decodeAuthLoginRequest(req)
		if err != nil {
			t.Errorf("error to decode from signInRequest function: %s", err.Error())
		}

		if dr.Email != payload.Email {
			t.Errorf("e-mail from decode request was incorrect, got: %s, want: %s", dr.Email, payload.Email)
		}

		if dr.Password != payload.Password {
			t.Errorf("password from decode request was incorrect, got: %s, want: %s", dr.Password, payload.Password)
		}
	}

}
