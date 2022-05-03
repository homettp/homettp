package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/homettp/homettp/internal/models"
	"github.com/petaki/support-go/forms"
)

func (a *app) commandIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.notFound(w)

		return
	}

	if r.Method != "GET" {
		a.methodNotAllowed(w, []string{"GET"})

		return
	}

	commands, err := a.commandRepository.FindAll()
	if err != nil {
		a.serverError(w, err)

		return
	}

	err = a.inertiaManager.Render(w, r, "command/Index", map[string]interface{}{
		"isCommandsActive": true,
		"commands":         commands,
	})
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) commandCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		a.getCommandCreate(w, r)
	case "POST":
		a.postCommandCreate(w, r)
	default:
		a.methodNotAllowed(w, []string{"GET", "POST"})
	}
}

func (a *app) getCommandCreate(w http.ResponseWriter, r *http.Request) {
	err := a.inertiaManager.Render(w, r, "command/Form", map[string]interface{}{
		"isCreateCommandActive": true,
		"command":               models.NewCommand(),
		"commandImages":         a.commandImages(),
		"commandPayload":        models.PayloadVariable,
		"errors":                forms.Bag{},
	})
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) postCommandCreate(w http.ResponseWriter, r *http.Request) {
	command := models.NewCommand()

	form, err := forms.NewFromRequest(w, r)
	if err != nil {
		a.formError(w, err)

		return
	}

	models.CommandCreateRules(form)

	if form.IsValid() {
		token, err := a.generateToken()
		if err != nil {
			a.serverError(w, err)
		}

		err = a.commandRepository.Create(command.Fill(form), token)
		if err != nil {
			switch err {
			case models.ErrDuplicateName:
				form.Errors.Add("name", "The name has already been taken.")
			case models.ErrInvalidValue:
				form.Errors.Add("value", "The value must contain command name.")
			default:
				a.serverError(w, err)

				return
			}
		} else {
			a.sessionManager.Put(r.Context(), sessionKeyFlashMessage, "Created successfully.")
			http.Redirect(w, r, fmt.Sprintf("/command/edit?id=%v", command.ID), http.StatusSeeOther)

			return
		}
	}

	err = a.inertiaManager.Render(w, r, "command/Form", map[string]interface{}{
		"isCreateCommandActive": true,
		"command":               command,
		"commandImages":         a.commandImages(),
		"commandPayload":        models.PayloadVariable,
		"errors":                form.Errors,
	})
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) commandEdit(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		a.getCommandEdit(w, r)
	case "POST":
		a.postCommandEdit(w, r)
	default:
		a.methodNotAllowed(w, []string{"GET", "POST"})
	}
}

func (a *app) getCommandEdit(w http.ResponseWriter, r *http.Request) {
	command, err := a.commandFromRequest(r, "id")
	if err != nil {
		a.notFound(w)

		return
	}

	err = a.inertiaManager.Render(w, r, "command/Form", map[string]interface{}{
		"command":        command,
		"commandImages":  a.commandImages(),
		"commandPayload": models.PayloadVariable,
		"commandPath":    command.Path(a.url),
		"errors":         forms.Bag{},
	})
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) postCommandEdit(w http.ResponseWriter, r *http.Request) {
	command, err := a.commandFromRequest(r, "id")
	if err != nil {
		a.notFound(w)

		return
	}

	form, err := forms.NewFromRequest(w, r)
	if err != nil {
		a.formError(w, err)

		return
	}

	models.CommandUpdateRules(form)

	if form.IsValid() {
		err = a.commandRepository.Update(command, (models.NewCommand()).Fill(form))
		if err != nil {
			switch err {
			case models.ErrDuplicateName:
				form.Errors.Add("name", "The name has already been taken.")
			case models.ErrInvalidValue:
				form.Errors.Add("value", "The value must contain command name.")
			default:
				a.serverError(w, err)

				return
			}
		} else {
			a.sessionManager.Put(r.Context(), sessionKeyFlashMessage, "Updated successfully.")
			http.Redirect(w, r, fmt.Sprintf("/command/edit?id=%v", command.ID), http.StatusSeeOther)

			return
		}
	}

	err = a.inertiaManager.Render(w, r, "command/Form", map[string]interface{}{
		"command":        command,
		"commandImages":  a.commandImages(),
		"commandPayload": models.PayloadVariable,
		"commandPath":    command.Path(a.url),
		"errors":         form.Errors,
	})
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) commandRefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		a.methodNotAllowed(w, []string{"PUT"})

		return
	}

	command, err := a.commandFromRequest(r, "id")
	if err != nil {
		a.notFound(w)

		return
	}

	token, err := a.generateToken()
	if err != nil {
		a.serverError(w, err)

		return
	}

	err = a.commandRepository.UpdateToken(command, token)
	if err != nil {
		a.serverError(w, err)

		return
	}

	a.sessionManager.Put(r.Context(), sessionKeyFlashMessage, "Updated successfully.")
	http.Redirect(w, r, fmt.Sprintf("/command/edit?id=%v", command.ID), http.StatusSeeOther)
}

func (a *app) commandDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		a.methodNotAllowed(w, []string{"DELETE"})

		return
	}

	command, err := a.commandFromRequest(r, "id")
	if err != nil {
		a.notFound(w)

		return
	}

	err = a.commandRepository.Delete(command)
	if err != nil {
		a.serverError(w, err)

		return
	}

	a.sessionManager.Put(r.Context(), sessionKeyFlashMessage, "Deleted successfully.")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *app) commandFromRequest(r *http.Request, parameter string) (*models.Command, error) {
	if r.URL.Query().Get(parameter) == "" {
		return nil, fmt.Errorf("%s parameter not found", parameter)
	}

	id, err := strconv.Atoi(r.URL.Query().Get(parameter))
	if err != nil {
		return nil, err
	}

	command, err := a.commandRepository.Find(id)
	if err != nil {
		return nil, err
	}

	return command, nil
}

func (a *app) commandImages() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":  "Door",
			"value": models.Door,
		},
		{
			"name":  "Light",
			"value": models.Light,
		},
		{
			"name":  "Outlet",
			"value": models.Outlet,
		},
		{
			"name":  "Plug",
			"value": models.Plug,
		},
		{
			"name":  "Sensor",
			"value": models.Sensor,
		},
	}
}
