package middleware

import (
	"net/http"

	"github.com/maqsatto/Notes-API/internal/logger"
)

func Recovery(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					log.Error("panic recovered", err)

					w.Write([]byte(`{"error": "internal server error"}`))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
