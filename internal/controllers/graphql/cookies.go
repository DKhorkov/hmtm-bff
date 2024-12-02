package graphqlcontroller

import (
	"net/http"
	"time"

	"github.com/DKhorkov/hmtm-bff/internal/config"
)

const (
	accessTokenCookieName  = "accessToken"
	refreshTokenCookieName = "refreshToken"
)

func setCookie(
	writer http.ResponseWriter,
	name string,
	value string,
	cookieConfig config.CookieConfig,
) {
	http.SetCookie(
		writer,
		&http.Cookie{
			Name:     name,
			Value:    value,
			HttpOnly: cookieConfig.HTTPOnly,
			Path:     cookieConfig.Path,
			Domain:   cookieConfig.Domain,
			Expires:  time.Now().Add(cookieConfig.Expires),
			MaxAge:   cookieConfig.MaxAge,
			SameSite: cookieConfig.SameSite,
			Secure:   cookieConfig.Secure,
		},
	)
}
