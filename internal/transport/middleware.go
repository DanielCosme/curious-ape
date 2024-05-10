package transport

import (
	"context"
	"errors"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/danielcosme/curious-ape/internal/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func midSlogConfig(t *Transport) middleware.RequestLoggerConfig {
	return middleware.RequestLoggerConfig{
		LogMethod:   true,
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: false, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				t.App.Log.LogAttrs(context.Background(), slog.LevelInfo, v.Method,
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				t.App.Log.LogAttrs(context.Background(), slog.LevelError, v.Method,
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
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
				c.Response().Header().Set("Vary", "Cookie")
				c.Response().Header().Set("Cache-Control", `no-cache="Set-Cookie"`)
			}
		})

		return next(c)
	}
}

func (t *Transport) midAuthenticateFromSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := t.SessionManager.GetInt(c.Request().Context(), ctxKeyAuthenticatedUserID)
		if id == 0 {
			return next(c)
		}
		exists, err := t.App.Exists(id)
		if err != nil {
			return err
		}
		if exists {
			ctx := context.WithValue(c.Request().Context(), ctxKeyIsAuthenticated, true)
			c.SetRequest(c.Request().WithContext(ctx))
		}
		return next(c)
	}
}

func (t *Transport) midRequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !t.IsAuthenticated(c.Request()) {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		// Set the "Cache-Control: no-store" header so that pages require
		// authentication are not stored in the users browser cache (or
		// other intermediary cache).
		c.Response().Header().Add("Cache-Control", "no-store")
		return next(c)
	}
}

func (t *Transport) midSetHabit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id < 1 {
			t.App.Log.Error("Habit id is invalid", "id", id)
			return errNotFound()
		}

		habit, err := t.App.HabitGetByID(id)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				return errNotFound()
			}
			return errServer(err)
		}

		c.Set("habit", habit)
		return next(c)
	}
}

func midSecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO make these headers appear ony in prod environments.
		// w.Header().Set("Content-Security-Policy",
		// 	"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		// w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		// w.Header().Set("X-Content-Type-Options", "nosniff")
		// w.Header().Set("X-Frame-Options", "deny")
		// w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (t *Transport) midRecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				t.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
