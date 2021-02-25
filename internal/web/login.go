package web

import (
	"errors"
	"net/http"

	"github.com/homettp/homettp/internal/models"
	"github.com/petaki/support-go/forms"
)

func (app *app) login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		app.getLogin(w, r)
	case "POST":
		app.postLogin(w, r)
	default:
		app.methodNotAllowed(w, []string{"GET", "POST"})
	}
}

func (app *app) getLogin(w http.ResponseWriter, r *http.Request) {
	err := app.inertiaManager.Render(w, r, "auth/Login", map[string]interface{}{
		"errors": forms.Bag{},
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *app) postLogin(w http.ResponseWriter, r *http.Request) {
	form, err := forms.NewFromRequest(w, r)
	if err != nil {
		app.formError(w, err)

		return
	}

	form.Required("username_or_email", "password")

	if form.IsValid() {
		user, err := app.userRepository.Authenticate(form.Data["username_or_email"].(string), form.Data["password"].(string))
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				form.Errors.Add("username_or_email", "These credentials do not match our records.")
			} else {
				app.serverError(w, err)

				return
			}
		} else {
			app.sessionManager.Put(r.Context(), sessionKeyAuthUserID, user.ID)

			err = app.sessionManager.RenewToken(r.Context())
			if err != nil {
				app.serverError(w, err)

				return
			}

			if form.Data["remember"].(bool) {
				if user.RememberToken == "" {
					token, err := app.generateToken()
					if err != nil {
						app.serverError(w, err)

						return
					}

					err = app.userRepository.UpdateRememberToken(user, token)
					if err != nil {
						app.serverError(w, err)

						return
					}
				}

				err = app.rememberCookie.SetValue(w, user.RememberCookie())
				if err != nil {
					app.serverError(w, err)

					return
				}
			}

			intendedURL := app.sessionManager.PopString(r.Context(), sessionKeyIntendedURL)
			if intendedURL != "" {
				http.Redirect(w, r, intendedURL, http.StatusSeeOther)

				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)

			return
		}
	}

	err = app.inertiaManager.Render(w, r, "auth/Login", map[string]interface{}{
		"errors": form.Errors,
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *app) logout(w http.ResponseWriter, r *http.Request) {
	app.sessionManager.Remove(r.Context(), sessionKeyAuthUserID)

	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)

		return
	}

	err = app.rememberCookie.Delete(w)
	if err != nil {
		app.serverError(w, err)

		return
	}

	user := app.authUser(r)

	if user != nil && user.RememberToken != "" {
		token, err := app.generateToken()
		if err != nil {
			app.serverError(w, err)

			return
		}

		err = app.userRepository.UpdateRememberToken(user, token)
		if err != nil {
			app.serverError(w, err)

			return
		}
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
