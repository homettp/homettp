package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/homettp/homettp/internal/models"
)

func (app *App) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *App) redirectIfNotAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.authUser(r) == nil {
			app.sessionManager.Put(r.Context(), sessionKeyIntendedUrl, r.URL.Path)
			http.Redirect(w, r, "/login", http.StatusFound)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *App) redirectIfAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.authUser(r) != nil {
			http.Redirect(w, r, "/", http.StatusFound)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *App) remember(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exists := app.sessionManager.Exists(r.Context(), sessionKeyAuthUserId)
		if exists {
			next.ServeHTTP(w, r)

			return
		}

		bytes, err := app.rememberCookie.GetValue(nil, r)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				next.ServeHTTP(w, r)

				return
			}

			app.serverError(w, err)

			return
		}

		token := string(bytes)

		if !strings.Contains(token, "|") {
			next.ServeHTTP(w, r)

			return
		}

		segments := strings.SplitN(token, "|", 2)

		if len(segments) != 2 || segments[0] == "" || segments[1] == "" {
			next.ServeHTTP(w, r)

			return
		}

		id, err := strconv.Atoi(segments[0])
		if err != nil {
			app.serverError(w, err)

			return
		}

		_, err = app.userRepository.AuthenticateByRememberToken(id, segments[1])
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				next.ServeHTTP(w, r)

				return
			}

			app.serverError(w, err)

			return
		}

		app.sessionManager.Put(r.Context(), sessionKeyAuthUserId, id)

		err = app.sessionManager.RenewToken(r.Context())
		if err != nil {
			app.serverError(w, err)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *App) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exists := app.sessionManager.Exists(r.Context(), sessionKeyAuthUserId)
		if !exists {
			next.ServeHTTP(w, r)

			return
		}

		user, err := app.userRepository.Find(app.sessionManager.GetInt(r.Context(), sessionKeyAuthUserId))
		if errors.Is(err, models.ErrNoRecord) || !user.IsEnabled {
			app.sessionManager.Remove(r.Context(), sessionKeyAuthUserId)
			next.ServeHTTP(w, r)

			return
		} else if err != nil {
			app.serverError(w, err)

			return
		}

		ctx := context.WithValue(r.Context(), contextKeyAuthUser, user)
		ctx = app.inertiaManager.WithProp(ctx, "auth", map[string]interface{}{
			"user":     user,
			"gravatar": user.Gravatar(88),
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *App) flashMessage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flashMessage := app.sessionManager.PopString(r.Context(), sessionKeyFlashMessage)
		if flashMessage == "" {
			next.ServeHTTP(w, r)

			return
		}

		ctx := app.inertiaManager.WithProp(r.Context(), "flash", flashMessage)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
