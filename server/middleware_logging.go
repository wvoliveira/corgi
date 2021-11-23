package server

import (
	"bytes"
	"context"
	"math/rand"
	"net/http"
	"time"
)

type LogResponseWriter struct {
	http.ResponseWriter
	statusCode int
	buf        bytes.Buffer
}

func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{ResponseWriter: w}
}

func (w *LogResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *LogResponseWriter) Write(body []byte) (int, error) {
	w.buf.Write(body)
	return w.ResponseWriter.Write(body)
}

func (Middlewares) getClientIP(r *http.Request) (ip string) {
	ip = r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	return
}

// AccessControl set common headers for web UI.
func (m Middlewares) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "TransactionId", rand.Uint64())
		r = r.WithContext(ctx)

		rw := NewLogResponseWriter(w)

		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		transactionId := r.Context().Value("TransactionId")

		m.logger.Infow("request",
			"transaction-id", transactionId,
			"account-id", r.Header.Get("AccountID"),
			"account-email", r.Header.Get("AccountEmail"),
			"account-role", r.Header.Get("AccountRole"),
			"host", r.Host,
			"remote-addr", m.getClientIP(r),
			"method", r.Method,
			"request-uri", r.RequestURI,
			"proto", r.Proto,
			"status", rw.statusCode,
			"content-len", len(rw.buf.String()),
			"user-agent", r.Header.Get("User-Agent"),
			"duration", duration,
		)

	})
}
