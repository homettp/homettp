package web

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/homettp/homettp/internal/models"
	"github.com/petaki/support-go/forms"
)

func (a *app) userIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		a.methodNotAllowed(w, []string{"GET", "HEAD"})

		return
	}

	users, err := a.userRepository.FindAll()
	if err != nil {
		a.serverError(w, err)

		return
	}

	gravatars := make(map[string]string, len(users))

	for _, user := range users {
		gravatars[user.Username] = user.Gravatar(96)
	}

	err = a.inertiaManager.Render(w, r, "user/Index", map[string]any{
		"isUsersActive": true,
		"users":         users,
		"gravatars":     gravatars,
	})
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) userCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET", "HEAD":
		a.getUserCreate(w, r)
	case "POST":
		a.postUserCreate(w, r)
	default:
		a.methodNotAllowed(w, []string{"GET", "HEAD", "POST"})
	}
}

func (a *app) getUserCreate(w http.ResponseWriter, r *http.Request) {
	err := a.inertiaManager.Render(w, r, "user/Form", map[string]any{
		"isCreateUserActive": true,
		"user":               models.NewUser(),
		"errors":             forms.Bag{},
	})
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) postUserCreate(w http.ResponseWriter, r *http.Request) {
	user := models.NewUser()

	form, err := forms.NewFromRequest(w, r)
	if err != nil {
		a.formError(w, err)

		return
	}

	models.UserCreateRules(form)

	if form.IsValid() {
		err = a.userRepository.Create(user.Fill(form))
		if err != nil {
			switch {
			case errors.Is(err, models.ErrDuplicateUsername):
				form.Errors.Add("username", "The username has already been taken.")
			case errors.Is(err, models.ErrDuplicateEmail):
				form.Errors.Add("email", "The email has already been taken.")
			default:
				a.serverError(w, err)

				return
			}
		} else {
			a.sessionManager.Put(r.Context(), sessionKeyFlashMessage, "Created successfully.")
			http.Redirect(w, r, fmt.Sprintf("/user/edit?id=%v", user.ID), http.StatusSeeOther)

			return
		}
	}

	err = a.inertiaManager.Render(w, r, "user/Form", map[string]any{
		"isCreateUserActive": true,
		"user":               user,
		"errors":             form.Errors,
	})
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) userEdit(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET", "HEAD":
		a.getUserEdit(w, r)
	case "POST":
		a.postUserEdit(w, r)
	default:
		a.methodNotAllowed(w, []string{"GET", "HEAD", "POST"})
	}
}

func (a *app) getUserEdit(w http.ResponseWriter, r *http.Request) {
	user, err := a.userFromRequest(r, "id")
	if err != nil {
		a.lookupError(w, err)

		return
	}

	err = a.inertiaManager.Render(w, r, "user/Form", map[string]any{
		"user":   user,
		"errors": forms.Bag{},
	})
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) postUserEdit(w http.ResponseWriter, r *http.Request) {
	user, err := a.userFromRequest(r, "id")
	if err != nil {
		a.lookupError(w, err)

		return
	}

	form, err := forms.NewFromRequest(w, r)
	if err != nil {
		a.formError(w, err)

		return
	}

	models.UserUpdateRules(form)

	if form.IsValid() {
		err = a.userRepository.Update(user, (models.NewUser()).Fill(form))
		if err != nil {
			switch {
			case errors.Is(err, models.ErrNoRecord):
				a.notFound(w)

				return
			case errors.Is(err, models.ErrDuplicateUsername):
				form.Errors.Add("username", "The username has already been taken.")
			case errors.Is(err, models.ErrDuplicateEmail):
				form.Errors.Add("email", "The email has already been taken.")
			default:
				a.serverError(w, err)

				return
			}
		} else {
			a.sessionManager.Put(r.Context(), sessionKeyFlashMessage, "Updated successfully.")
			http.Redirect(w, r, fmt.Sprintf("/user/edit?id=%v", user.ID), http.StatusSeeOther)

			return
		}
	}

	err = a.inertiaManager.Render(w, r, "user/Form", map[string]any{
		"user":   user,
		"errors": form.Errors,
	})
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) userDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		a.methodNotAllowed(w, []string{"DELETE"})

		return
	}

	user, err := a.userFromRequest(r, "id")
	if err != nil {
		a.lookupError(w, err)

		return
	}

	err = a.userRepository.Delete(user)
	if errors.Is(err, models.ErrNoRecord) {
		a.notFound(w)

		return
	}

	if err != nil {
		a.serverError(w, err)

		return
	}

	a.sessionManager.Put(r.Context(), sessionKeyFlashMessage, "Deleted successfully.")
	http.Redirect(w, r, "/user", http.StatusSeeOther)
}

func (a *app) userFromRequest(r *http.Request, parameter string) (*models.User, error) {
	if r.URL.Query().Get(parameter) == "" {
		return nil, fmt.Errorf("%w: %s", ErrMissingParameter, parameter)
	}

	id, err := strconv.Atoi(r.URL.Query().Get(parameter))
	if err != nil {
		return nil, err
	}

	user, err := a.userRepository.Find(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
