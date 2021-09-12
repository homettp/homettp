package web

import (
	"errors"
	"net/http"

	"github.com/homettp/homettp/internal/models"
	"github.com/petaki/support-go/forms"
)

func (a *app) login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		a.getLogin(w, r)
	case "POST":
		a.postLogin(w, r)
	default:
		a.methodNotAllowed(w, []string{"GET", "POST"})
	}
}

func (a *app) getLogin(w http.ResponseWriter, r *http.Request) {
	err := a.inertiaManager.Render(w, r, "auth/Login", map[string]interface{}{
		"errors": forms.Bag{},
	})
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) postLogin(w http.ResponseWriter, r *http.Request) {
	form, err := forms.NewFromRequest(w, r)
	if err != nil {
		a.formError(w, err)

		return
	}

	form.Required("username_or_email", "password")

	if form.IsValid() {
		user, err := a.userRepository.Authenticate(form.Data["username_or_email"].(string), form.Data["password"].(string))
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				form.Errors.Add("username_or_email", "These credentials do not match our records.")
			} else {
				a.serverError(w, err)

				return
			}
		} else {
			a.sessionManager.Put(r.Context(), sessionKeyAuthUserID, user.ID)

			err = a.sessionManager.RenewToken(r.Context())
			if err != nil {
				a.serverError(w, err)

				return
			}

			if form.Data["remember"].(bool) {
				if user.RememberToken == "" {
					token, err := a.generateToken()
					if err != nil {
						a.serverError(w, err)

						return
					}

					err = a.userRepository.UpdateRememberToken(user, token)
					if err != nil {
						a.serverError(w, err)

						return
					}
				}

				err = a.rememberCookie.SetValue(w, user.RememberCookie())
				if err != nil {
					a.serverError(w, err)

					return
				}
			}

			intendedURL := a.sessionManager.PopString(r.Context(), sessionKeyIntendedURL)
			if intendedURL != "" {
				http.Redirect(w, r, intendedURL, http.StatusSeeOther)

				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)

			return
		}
	}

	err = a.inertiaManager.Render(w, r, "auth/Login", map[string]interface{}{
		"errors": form.Errors,
	})
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) logout(w http.ResponseWriter, r *http.Request) {
	a.sessionManager.Remove(r.Context(), sessionKeyAuthUserID)

	err := a.sessionManager.RenewToken(r.Context())
	if err != nil {
		a.serverError(w, err)

		return
	}

	err = a.rememberCookie.Delete(w)
	if err != nil {
		a.serverError(w, err)

		return
	}

	user := a.authUser(r)

	if user != nil && user.RememberToken != "" {
		token, err := a.generateToken()
		if err != nil {
			a.serverError(w, err)

			return
		}

		err = a.userRepository.UpdateRememberToken(user, token)
		if err != nil {
			a.serverError(w, err)

			return
		}
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
