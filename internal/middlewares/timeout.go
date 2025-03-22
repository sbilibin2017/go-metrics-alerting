package middlewares

import (
	"context"
	"net/http"
	"time"
)

func TimeoutMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), time.Duration(5*time.Second))
			defer cancel()

			r = r.WithContext(ctx)

			done := make(chan struct{})

			go func() {
				next.ServeHTTP(w, r)
				close(done)
			}()

			select {
			case <-ctx.Done():
				http.Error(w, "Request Timeout", http.StatusRequestTimeout)
			case <-done:
			}
		})
	}
}
