package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// TimeoutMiddleware создает middleware для установки таймаута на каждый запрос.
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r = r.WithContext(ctx)

			done := make(chan struct{})

			go func() {
				next.ServeHTTP(w, r)
				close(done)
			}()

			select {
			case <-done:
			case <-ctx.Done():
				http.Error(w, fmt.Sprintf("Request timed out after %v", timeout), http.StatusRequestTimeout)
			}
		})
	}
}
