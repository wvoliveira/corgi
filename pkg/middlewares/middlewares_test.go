package middlewares

import (
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccess(t *testing.T) {
	logger, entries := log.NewForTest()
	handler := Access(logger)

	router := gin.New()
	router.Use(handler)

	// RUN
	_ = performRequest(router, "GET", "/")

	// TEST
	assert.Equal(t, 2, entries.Len())
	assert.Equal(t, "", entries.All()[0].Message)
}

type header struct {
	Key   string
	Value string
}

func performRequest(r http.Handler, method, path string, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
