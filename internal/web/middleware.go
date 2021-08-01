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

func (a *app) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				a.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (a *app) redirectIfNotAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.authUser(r) == nil {
			a.sessionManager.Put(r.Context(), sessionKeyIntendedURL, r.URL.Path)
			http.Redirect(w, r, "/login", http.StatusFound)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *app) redirectIfAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.authUser(r) != nil {
			http.Redirect(w, r, "/", http.StatusFound)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *app) remember(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exists := a.sessionManager.Exists(r.Context(), sessionKeyAuthUserID)
		if exists {
			next.ServeHTTP(w, r)

			return
		}

		bytes, err := a.rememberCookie.GetValue(nil, r)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				next.ServeHTTP(w, r)

				return
			}

			a.serverError(w, err)

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
			a.serverError(w, err)

			return
		}

		_, err = a.userRepository.AuthenticateByRememberToken(id, segments[1])
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				next.ServeHTTP(w, r)

				return
			}

			a.serverError(w, err)

			return
		}

		a.sessionManager.Put(r.Context(), sessionKeyAuthUserID, id)

		err = a.sessionManager.RenewToken(r.Context())
		if err != nil {
			a.serverError(w, err)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *app) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exists := a.sessionManager.Exists(r.Context(), sessionKeyAuthUserID)
		if !exists {
			next.ServeHTTP(w, r)

			return
		}

		user, err := a.userRepository.Find(a.sessionManager.GetInt(r.Context(), sessionKeyAuthUserID))
		if errors.Is(err, models.ErrNoRecord) || !user.IsEnabled {
			a.sessionManager.Remove(r.Context(), sessionKeyAuthUserID)
			next.ServeHTTP(w, r)

			return
		} else if err != nil {
			a.serverError(w, err)

			return
		}

		ctx := context.WithValue(r.Context(), contextKeyAuthUser, user)
		ctx = a.inertiaManager.WithProp(ctx, "auth", map[string]interface{}{
			"user":     user,
			"gravatar": user.Gravatar(88),
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *app) flashMessage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flashMessage := a.sessionManager.PopString(r.Context(), sessionKeyFlashMessage)
		if flashMessage == "" {
			next.ServeHTTP(w, r)

			return
		}

		ctx := a.inertiaManager.WithProp(r.Context(), "flash", flashMessage)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
