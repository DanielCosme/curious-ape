package api

import (
	"context"
	"time"

	"git.danicos.dev/daniel/curious-ape/pkg/dove"
	"git.danicos.dev/daniel/curious-ape/pkg/oak"
	"github.com/alexedwards/scs/v2"
)

func (a *API) MiddlewareLoadCookie(next dove.HandlerFunc) dove.HandlerFunc {
	return func(c *dove.Context) error {
		c.Res.Header().Add("Vary", "Cookie")
		var token string
		cookie, err := c.Req.Cookie(a.Scs.Cookie.Name)
		if err == nil {
			token = cookie.Value
		}
		ctx, err := a.Scs.Load(c.Ctx(), token)
		if err != nil {
			return err
		}
		c.Req = c.Req.WithContext(ctx)

		c.Res.Before(func() {
			switch a.Scs.Status(ctx) {
			case scs.Modified:
				token, expiry, err := a.Scs.Commit(ctx)
				if err != nil {
					panic(err)
				}
				a.Scs.WriteSessionCookie(ctx, c.Res.Writer, token, expiry)
			case scs.Destroyed:
				a.Scs.WriteSessionCookie(ctx, c.Res.Writer, "", time.Time{})
			}
		})
		return next(c)
	}
}

func (a *API) MiddlewareAuthenticateFromSession(next dove.HandlerFunc) dove.HandlerFunc {
	return func(c *dove.Context) error {
		logger := oak.FromContext(c.Ctx())
		id := a.Scs.GetInt(c.Ctx(), string(ctxKeyAuthenticatedUserID))
		if id == 0 {
			logger.Warning("no authenticated user found in session")
			return next(c)
		}
		usr, err := a.App.GetUser(id)
		if err != nil {
			return err
		}
		logger.Debug("authenticated from session", "username", usr.Username)
		ctx := context.WithValue(c.Ctx(), ctxKeyIsAuthenticated, true)
		ctx = context.WithValue(ctx, ctxUser, usr)
		c.Req = c.Req.WithContext(ctx)
		return next(c)
	}
}

func (a *API) MiddlewareRequireAuthentication(next dove.HandlerFunc) dove.HandlerFunc {
	return func(c *dove.Context) error {
		if !a.IsAuthenticated(c.Req) {
			return c.Redirect("/login")
		}
		// Set the "Cache-Control: no-store" header so that pages require
		// authentication are not stored in the users browser cache (or
		// other intermediary cache).
		c.Res.Header().Add("Cache-Control", "no-store")
		return next(c)
	}
}

func DevMiddleware(next dove.HandlerFunc) dove.HandlerFunc {
	return func(c *dove.Context) error {
		// Don't cache in development.
		c.Res.Header().Add("Cache-Control", "no-store")
		return next(c)
	}
}
