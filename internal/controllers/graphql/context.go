package graphqlcontroller

import (
	"context"
	"net/http"

	"github.com/DKhorkov/hmtm-bff/internal/middlewares"
)

func getCookieFromContext(ctx context.Context, cookieName string) (*http.Cookie, bool) {
	cookie, ok := ctx.Value(middlewares.ContextKey(cookieName)).(*http.Cookie)
	return cookie, ok
}

func getHTTPWriterFromContext(ctx context.Context) (http.ResponseWriter, bool) {
	writer, ok := ctx.Value(middlewares.ContextKey(middlewares.CookiesWriterName)).(http.ResponseWriter)
	return writer, ok
}
