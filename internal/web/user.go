package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/homettp/homettp/internal/models"
	"github.com/petaki/support-go/forms"
)

func (a *app) userIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		a.methodNotAllowed(w, []string{"GET"})

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

	err = a.inertiaManager.Render(w, r, "user/Index", map[string]interface{}{
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
	case "GET":
		a.getUserCreate(w, r)
	case "POST":
		a.postUserCreate(w, r)
	default:
		a.methodNotAllowed(w, []string{"GET", "POST"})
	}
}

func (a *app) getUserCreate(w http.ResponseWriter, r *http.Request) {
	err := a.inertiaManager.Render(w, r, "user/Form", map[string]interface{}{
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
			switch err {
			case models.ErrDuplicateUsername:
				form.Errors.Add("username", "The username has already been taken.")
			case models.ErrDuplicateEmail:
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

	err = a.inertiaManager.Render(w, r, "user/Form", map[string]interface{}{
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
	case "GET":
		a.getUserEdit(w, r)
	case "POST":
		a.postUserEdit(w, r)
	default:
		a.methodNotAllowed(w, []string{"GET", "POST"})
	}
}

func (a *app) getUserEdit(w http.ResponseWriter, r *http.Request) {
	user, err := a.userFromRequest(r, "id")
	if err != nil {
		a.notFound(w)

		return
	}

	err = a.inertiaManager.Render(w, r, "user/Form", map[string]interface{}{
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
		a.notFound(w)

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
			switch err {
			case models.ErrDuplicateUsername:
				form.Errors.Add("username", "The username has already been taken.")
			case models.ErrDuplicateEmail:
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

	err = a.inertiaManager.Render(w, r, "user/Form", map[string]interface{}{
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
		a.notFound(w)

		return
	}

	err = a.userRepository.Delete(user)
	if err != nil {
		a.serverError(w, err)

		return
	}

	a.sessionManager.Put(r.Context(), sessionKeyFlashMessage, "Deleted successfully.")
	http.Redirect(w, r, "/user", http.StatusSeeOther)
}

func (a *app) userFromRequest(r *http.Request, parameter string) (*models.User, error) {
	if r.URL.Query().Get(parameter) == "" {
		return nil, fmt.Errorf("%s parameter not found", parameter)
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
