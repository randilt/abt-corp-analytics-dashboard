package middleware

import (
	"net/http"
	"runtime/debug"

	"analytics-dashboard-api/internal/utils"
	"analytics-dashboard-api/pkg/logger"
)

// Recovery middleware for panic recovery
func Recovery(logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Panic recovered",
						"error", err,
						"path", r.URL.Path,
						"method", r.Method,
						"stack", string(debug.Stack()),
					)

					utils.WriteErrorResponse(w, http.StatusInternalServerError, 
						"Internal server error occurred")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}