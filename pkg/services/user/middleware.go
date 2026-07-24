package user

import (
	"context"
	"log/slog"
	"net/http"

	"danicos.dev/daniel/curious-ape/pkg/web"
	"github.com/alexedwards/scs/v2"
	"github.com/stephenafamo/bob"
)

func MiddlewareAuthenticateFromSession(session *scs.SessionManager, db bob.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := session.GetInt64(r.Context(), string(ctxKeyAuthenticatedUserID))
			if id == 0 {
				next.ServeHTTP(w, r)
				return
			}

			user, err := Get(db, Params{ID: id})
			if err != nil {
				web.ErrInternalServer(err, w)
				return
			}
			slog.Debug("User authenticated from session", "username", user.Name)

			ctx := context.WithValue(r.Context(), ctxKeyIsAuthenticated, true)
			ctx = context.WithValue(ctx, ctxUser, user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
