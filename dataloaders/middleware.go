package dataloaders

import (
	"context"
	"github.com/RobusK/GoCtaApi/api"
	"net/http"
)

// Middleware stores Loaders as a request-scoped context value.
func Middleware(client api.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			loaders := newLoaders(ctx, client)
			augmentedCtx := context.WithValue(ctx, key, loaders)
			r = r.WithContext(augmentedCtx)
			next.ServeHTTP(w, r)
		})
	}
}
