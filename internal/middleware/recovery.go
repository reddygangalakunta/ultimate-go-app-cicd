package middleware

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/example/ultimate-ci-cd-pipeline/internal/model"
)

// PanicRecovery catches panics and returns a 500 error instead of crashing the server.
func PanicRecovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					stack := debug.Stack()
					logger.Error("Unhandled Panic Recovered",
						slog.Any("error", err),
						slog.String("stack", string(stack)),
						slog.String("path", r.URL.Path),
					)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(model.ErrorResponse{
						Code:    http.StatusInternalServerError,
						Message: fmt.Sprintf("Internal Server Error: %v", err),
					})
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
