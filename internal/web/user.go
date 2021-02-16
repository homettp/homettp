package web

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/homettp/homettp/internal/models"
	"github.com/petaki/support-go/forms"
)

func (app *App) userIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		app.methodNotAllowed(w, []string{"GET"})

		return
	}

	users, err := app.userRepository.FindAll()
	if err != nil {
		app.serverError(w, err)

		return
	}

	gravatars := make(map[string]string, len(users))

	for _, user := range users {
		gravatars[user.Username] = user.Gravatar(96)
	}

	err = app.inertiaManager.Render(w, r, "user/Index", map[string]interface{}{
		"isUsersActive": true,
		"users":         users,
		"gravatars":     gravatars,
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *App) userCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		app.getUserCreate(w, r)
	case "POST":
		app.postUserCreate(w, r)
	default:
		app.methodNotAllowed(w, []string{"GET", "POST"})
	}
}

func (app *App) getUserCreate(w http.ResponseWriter, r *http.Request) {
	err := app.inertiaManager.Render(w, r, "user/Form", map[string]interface{}{
		"isCreateUserActive": true,
		"user":               models.NewUser(),
		"errors":             forms.Bag{},
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *App) postUserCreate(w http.ResponseWriter, r *http.Request) {
	user := models.NewUser()

	form, err := forms.NewFromRequest(w, r)
	if err != nil {
		app.formError(w, err)

		return
	}

	models.UserCreateRules(form)

	if form.IsValid() {
		err = app.userRepository.Create(user.Fill(form))
		if err != nil {
			switch err {
			case models.ErrDuplicateUsername:
				form.Errors.Add("username", "The username has already been taken.")
			case models.ErrDuplicateEmail:
				form.Errors.Add("email", "The email has already been taken.")
			default:
				app.serverError(w, err)

				return
			}
		} else {
			app.sessionManager.Put(r.Context(), sessionKeyFlashMessage, "Created successfully.")
			http.Redirect(w, r, fmt.Sprintf("/user/edit?id=%v", user.Id), http.StatusSeeOther)

			return
		}
	}

	err = app.inertiaManager.Render(w, r, "user/Form", map[string]interface{}{
		"isCreateUserActive": true,
		"user":               user,
		"errors":             form.Errors,
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *App) userEdit(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		app.getUserEdit(w, r)
	case "POST":
		app.postUserEdit(w, r)
	default:
		app.methodNotAllowed(w, []string{"GET", "POST"})
	}
}

func (app *App) getUserEdit(w http.ResponseWriter, r *http.Request) {
	user, err := app.userFromRequest(r, "id")
	if err != nil {
		app.notFound(w)

		return
	}

	err = app.inertiaManager.Render(w, r, "user/Form", map[string]interface{}{
		"user":   user,
		"errors": forms.Bag{},
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *App) postUserEdit(w http.ResponseWriter, r *http.Request) {
	user, err := app.userFromRequest(r, "id")
	if err != nil {
		app.notFound(w)

		return
	}

	form, err := forms.NewFromRequest(w, r)
	if err != nil {
		app.formError(w, err)

		return
	}

	models.UserUpdateRules(form)

	if form.IsValid() {
		err = app.userRepository.Update(user, (models.NewUser()).Fill(form))
		if err != nil {
			switch err {
			case models.ErrDuplicateUsername:
				form.Errors.Add("username", "The username has already been taken.")
			case models.ErrDuplicateEmail:
				form.Errors.Add("email", "The email has already been taken.")
			default:
				app.serverError(w, err)

				return
			}
		} else {
			app.sessionManager.Put(r.Context(), sessionKeyFlashMessage, "Updated successfully.")
			http.Redirect(w, r, fmt.Sprintf("/user/edit?id=%v", user.Id), http.StatusSeeOther)

			return
		}
	}

	err = app.inertiaManager.Render(w, r, "user/Form", map[string]interface{}{
		"user":   user,
		"errors": form.Errors,
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *App) userDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		app.methodNotAllowed(w, []string{"DELETE"})

		return
	}

	user, err := app.userFromRequest(r, "id")
	if err != nil {
		app.notFound(w)

		return
	}

	err = app.userRepository.Delete(user)
	if err != nil {
		app.serverError(w, err)

		return
	}

	app.sessionManager.Put(r.Context(), sessionKeyFlashMessage, "Deleted successfully.")
	http.Redirect(w, r, "/user", http.StatusSeeOther)
}

func (app *App) userFromRequest(r *http.Request, parameter string) (*models.User, error) {
	if r.URL.Query().Get(parameter) == "" {
		return nil, errors.New(fmt.Sprintf("%s parameter not found", parameter))
	}

	id, err := strconv.Atoi(r.URL.Query().Get(parameter))
	if err != nil {
		return nil, err
	}

	user, err := app.userRepository.Find(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
