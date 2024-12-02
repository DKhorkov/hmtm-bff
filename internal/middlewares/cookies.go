package middlewares

import (
	"context"
	"net/http"
)

var CookiesWriterName = "cookiesWriterName"

// ContextKey - context-keys-type: should not use basic type string as key in context.WithValue (revive)
// https://vishnubharathi.codes/blog/context-with-value-pitfall/
type ContextKey string

// CookiesMiddleware reads provided cookies from request and paste them into context for graphql purposes.
func CookiesMiddleware(handler http.Handler, cookieNames []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, cookieName := range cookieNames {
			cookie, err := r.Cookie(cookieName)
			if err != nil {
				continue
			}

			ctx := context.WithValue(r.Context(), ContextKey(cookieName), cookie)
			r = r.WithContext(ctx)
		}

		// Paste writer to context for writing cookies in resolvers purposes:
		ctx := context.WithValue(r.Context(), ContextKey(CookiesWriterName), w)
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	})
}
