package api

import (
	"context"
	"github.com/alexedwards/scs/v2"
	"github.com/danielcosme/curious-ape/pkg/application"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log/slog"
	"net/http"
	"time"
)

func midSlogConfig(t *Transport) middleware.RequestLoggerConfig {
	return middleware.RequestLoggerConfig{
		LogMethod:   true,
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: false, // forwards error to the global error handler, so it can decide appropriate status code
		LogRemoteIP: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				t.App.Log.LogAttrs(context.Background(), slog.LevelInfo, v.Method,
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("IP", v.RemoteIP),
				)
			} else {
				t.App.Log.LogAttrs(context.Background(), slog.LevelError, v.Method,
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("IP", v.RemoteIP),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}
}

func (t *Transport) midLoadAndSaveCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var token string
		cookie, err := c.Cookie(t.SessionManager.Cookie.Name)
		if err == nil {
			token = cookie.Value
		}

		ctx, err = t.SessionManager.Load(ctx, token)
		if err != nil {
			return err
		}

		c.SetRequest(c.Request().WithContext(ctx))

		c.Response().Before(func() {
			if t.SessionManager.Status(ctx) != scs.Unmodified {
				responseCookie := &http.Cookie{
					Name:     t.SessionManager.Cookie.Name,
					Path:     t.SessionManager.Cookie.Path,
					Domain:   t.SessionManager.Cookie.Domain,
					Secure:   t.SessionManager.Cookie.Secure,
					HttpOnly: t.SessionManager.Cookie.HttpOnly,
					SameSite: t.SessionManager.Cookie.SameSite,
				}

				switch t.SessionManager.Status(ctx) {
				case scs.Modified:
					token, _, err = t.SessionManager.Commit(ctx)
					if err != nil {
						panic(err)
					}

					responseCookie.Value = token
				case scs.Destroyed:
					responseCookie.Expires = time.Unix(1, 0)
					responseCookie.MaxAge = -1
				}

				c.SetCookie(responseCookie)
				c.Response().Header().Add("Vary", "Cookie")
				c.Response().Header().Add("Cache-Control", `no-cache="Set-Cookie"`)
			}
		})

		return next(c)
	}
}

func (t *Transport) midAuthenticateFromSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := t.SessionManager.GetInt(c.Request().Context(), string(ctxKeyAuthenticatedUserID))
		if id == 0 {
			return next(c)
		}
		usr, err := t.App.GetUser(id)
		if err != nil {
			return err
		}
		if usr != nil {
			ctx := context.WithValue(c.Request().Context(), ctxKeyIsAuthenticated, true)
			ctx = context.WithValue(ctx, ctxUser, usr)
			c.SetRequest(c.Request().WithContext(ctx))
		}
		return next(c)
	}
}

func (t *Transport) midRequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !t.IsAuthenticated(c.Request()) {
			return c.NoContent(http.StatusUnauthorized)
		}

		// Set the "Cache-Control: no-store" header so that pages require
		// authentication are not stored in the users browser cache (or
		// other intermediary cache).
		c.Response().Header().Add("Cache-Control", "no-store")
		return next(c)
	}
}

func (t *Transport) midSecureHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		if t.App.Env == application.Dev {
			// NOTE(daniel): To make vue/vite work on dev mode
			c.Response().Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			c.Response().Header().Set("Access-Control-Allow-Origin", "https://danicos.me")
		}
		c.Response().Header().Set("Referrer-Policy", "origin-when-cross-origin")

		// c.Response().Header().Set("X-Content-Kind-Options", "nosniff")
		// c.Response().Header().Set("X-Frame-Options", "deny")
		// c.Response().Header().Set("X-XSS-Protection", "0")
		return next(c)
	}
}
