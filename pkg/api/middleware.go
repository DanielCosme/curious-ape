package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func midSlogConfig(t *API) middleware.RequestLoggerConfig {
	return middleware.RequestLoggerConfig{
		LogMethod:   true,
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: false, // forwards error to the global error handler, so it can decide appropriate status code
		LogRemoteIP: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []slog.Attr{
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.String("IP", v.RemoteIP),
				slog.String("Duration", fmt.Sprintf("%d ms", time.Since(v.StartTime).Milliseconds())),
			}
			if v.Error == nil {
				t.App.Log.LogAttrs(ctx(c), slog.LevelInfo, v.Method, attrs...)
			} else {
				t.App.Log.LogAttrs(ctx(c), slog.LevelError, v.Method,
					append(attrs, slog.String("err", v.Error.Error()))...,
				)
			}
			return nil
		},
	}
}

func (api *API) midLoadAndSaveCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Add("Vary", "Cookie")
		ctx := c.Request().Context()

		var token string
		cookie, err := c.Cookie(api.SessionManager.Cookie.Name)
		if err == nil {
			token = cookie.Value
		}
		ctx, err = api.SessionManager.Load(ctx, token)
		if err != nil {
			return err
		}
		c.SetRequest(c.Request().WithContext(ctx))

		c.Response().Before(func() {
			switch api.SessionManager.Status(ctx) {
			case scs.Modified:
				token, expiry, err := api.SessionManager.Commit(ctx)
				if err != nil {
					panic(err)
				}
				api.SessionManager.WriteSessionCookie(ctx, c.Response().Writer, token, expiry)
			case scs.Destroyed:
				api.SessionManager.WriteSessionCookie(ctx, c.Response().Writer, "", time.Time{})
			}
		})

		return next(c)
	}
}

func (api *API) midAuthenticateFromSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := api.SessionManager.GetInt(c.Request().Context(), string(ctxKeyAuthenticatedUserID))
		if id == 0 {
			return next(c)
		}
		usr, err := api.App.GetUser(int(id))
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

func (api *API) midRequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !api.IsAuthenticated(c.Request()) {
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}

		// Set the "Cache-Control: no-store" header so that pages require
		// authentication are not stored in the users browser cache (or
		// other intermediary cache).
		c.Response().Header().Add("Cache-Control", "no-store")
		return next(c)
	}
}

func (api *API) midSecureHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: Figure out secure headers for good.
		// c.Response().Header().Set("Content-Security-Policy",
		// 	"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		// c.Response().Header().Set("Access-Control-Allow-Origin", "https://danicos.me")
		// c.Response().Header().Set("Referrer-Policy", "origin-when-cross-origin")

		// c.Response().Header().Set("X-Content-Kind-Options", "nosniff")
		// c.Response().Header().Set("X-Frame-Options", "deny")
		// c.Response().Header().Set("X-XSS-Protection", "0")
		return next(c)
	}
}

func ctx(c echo.Context) context.Context {
	return c.Request().Context()
}
