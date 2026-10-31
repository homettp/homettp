package web

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/homettp/homettp/internal/models"
	"github.com/petaki/inertia-go"
)

func (a *app) preventCSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET", "HEAD", "OPTIONS":
			next.ServeHTTP(w, r)

			return
		}

		if r.Header.Get(inertia.HeaderInertia) == "" || !a.isSameOrigin(r) {
			a.clientError(w, http.StatusForbidden)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *app) storeIntendedURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.authUser(r) == nil && r.Method == "GET" && isKnownPath(r) {
			a.sessionManager.Put(r.Context(), sessionKeyIntendedURL, r.URL.RequestURI())
		}

		next.ServeHTTP(w, r)
	})
}

func (a *app) redirectIfNotAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.authUser(r) == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *app) redirectIfAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.authUser(r) != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)

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

			a.forgetRemember(w, r, next)

			return
		}

		segments := strings.SplitN(string(bytes), "|", 2)

		if len(segments) != 2 || segments[0] == "" || segments[1] == "" {
			a.forgetRemember(w, r, next)

			return
		}

		id, err := strconv.Atoi(segments[0])
		if err != nil {
			a.forgetRemember(w, r, next)

			return
		}

		_, err = a.userRepository.AuthenticateByRememberToken(id, segments[1])
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				a.forgetRemember(w, r, next)

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

func (a *app) forgetRemember(w http.ResponseWriter, r *http.Request, next http.Handler) {
	err := a.rememberCookie.Delete(w)
	if err != nil {
		a.serverError(w, err)

		return
	}

	next.ServeHTTP(w, r)
}

func (a *app) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exists := a.sessionManager.Exists(r.Context(), sessionKeyAuthUserID)
		if !exists {
			next.ServeHTTP(w, r)

			return
		}

		user, err := a.userRepository.Find(a.sessionManager.GetInt(r.Context(), sessionKeyAuthUserID))
		if err != nil && !errors.Is(err, models.ErrNoRecord) {
			a.serverError(w, err)

			return
		}

		if err != nil || !user.IsEnabled {
			a.sessionManager.Remove(r.Context(), sessionKeyAuthUserID)
			next.ServeHTTP(w, r)

			return
		}

		ctx := context.WithValue(r.Context(), contextKeyAuthUser, user)
		ctx = a.inertiaManager.WithProp(ctx, "auth", map[string]any{
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

		ctx := a.inertiaManager.WithFlash(r.Context(), map[string]any{
			"success": flashMessage,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
