package web

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/homettp/homettp/internal/models"
	"github.com/petaki/support-go/forms"
)

func (app *App) commandIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)

		return
	}

	if r.Method != "GET" {
		app.methodNotAllowed(w, []string{"GET"})

		return
	}

	commands, err := app.commandRepository.FindAll()
	if err != nil {
		app.serverError(w, err)

		return
	}

	err = app.inertiaManager.Render(w, r, "command/Index", map[string]interface{}{
		"isCommandsActive": true,
		"commands":         commands,
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *App) commandCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		app.getCommandCreate(w, r)
	case "POST":
		app.postCommandCreate(w, r)
	default:
		app.methodNotAllowed(w, []string{"GET", "POST"})
	}
}

func (app *App) getCommandCreate(w http.ResponseWriter, r *http.Request) {
	err := app.inertiaManager.Render(w, r, "command/Form", map[string]interface{}{
		"isCreateCommandActive": true,
		"command":               models.NewCommand(),
		"commandImages":         app.getCommandImages(),
		"errors":                forms.Bag{},
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *App) postCommandCreate(w http.ResponseWriter, r *http.Request) {
	command := models.NewCommand()

	form, err := forms.NewFromRequest(w, r)
	if err != nil {
		app.formError(w, err)

		return
	}

	models.CommandCreateRules(form)

	if form.IsValid() {
		token, err := app.generateToken()
		if err != nil {
			app.serverError(w, err)
		}

		err = app.commandRepository.Create(command.Fill(form), token)
		if err != nil {
			switch err {
			case models.ErrDuplicateName:
				form.Errors.Add("name", "The name has already been taken.")
			default:
				app.serverError(w, err)

				return
			}
		} else {
			app.sessionManager.Put(r.Context(), sessionKeyFlashMessage, "Created successfully.")
			http.Redirect(w, r, "/", http.StatusSeeOther)

			return
		}
	}

	err = app.inertiaManager.Render(w, r, "command/Form", map[string]interface{}{
		"isCreateCommandActive": true,
		"command":               command,
		"commandImages":         app.getCommandImages(),
		"errors":                form.Errors,
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *App) commandEdit(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		app.getCommandEdit(w, r)
	case "POST":
		app.postCommandEdit(w, r)
	default:
		app.methodNotAllowed(w, []string{"GET", "POST"})
	}
}

func (app *App) getCommandEdit(w http.ResponseWriter, r *http.Request) {
	command, err := app.getCommandFromRequest(r, "id")
	if err != nil {
		app.notFound(w)

		return
	}

	err = app.inertiaManager.Render(w, r, "command/Form", map[string]interface{}{
		"command":       command,
		"commandImages": app.getCommandImages(),
		"commandPath":   command.Path(app.url),
		"errors":        forms.Bag{},
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *App) postCommandEdit(w http.ResponseWriter, r *http.Request) {
	command, err := app.getCommandFromRequest(r, "id")
	if err != nil {
		app.notFound(w)

		return
	}

	form, err := forms.NewFromRequest(w, r)
	if err != nil {
		app.formError(w, err)

		return
	}

	models.CommandUpdateRules(form)

	if form.IsValid() {
		err = app.commandRepository.Update(command, (models.NewCommand()).Fill(form))
		if err != nil {
			switch err {
			case models.ErrDuplicateName:
				form.Errors.Add("name", "The name has already been taken.")
			default:
				app.serverError(w, err)

				return
			}
		} else {
			app.sessionManager.Put(r.Context(), sessionKeyFlashMessage, "Updated successfully.")
			http.Redirect(w, r, "/", http.StatusSeeOther)

			return
		}
	}

	err = app.inertiaManager.Render(w, r, "command/Form", map[string]interface{}{
		"command":       command,
		"commandImages": app.getCommandImages(),
		"commandPath":   command.Path(app.url),
		"errors":        form.Errors,
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *App) commandRefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		app.methodNotAllowed(w, []string{"PUT"})

		return
	}

	command, err := app.getCommandFromRequest(r, "id")
	if err != nil {
		app.notFound(w)

		return
	}

	token, err := app.generateToken()
	if err != nil {
		app.notFound(w)

		return
	}

	err = app.commandRepository.UpdateToken(command, token)
	if err != nil {
		app.serverError(w, err)

		return
	}
}

func (app *App) commandDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		app.methodNotAllowed(w, []string{"DELETE"})

		return
	}

	command, err := app.getCommandFromRequest(r, "id")
	if err != nil {
		app.notFound(w)

		return
	}

	err = app.commandRepository.Delete(command)
	if err != nil {
		app.serverError(w, err)

		return
	}

	app.sessionManager.Put(r.Context(), sessionKeyFlashMessage, "Deleted successfully.")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *App) getCommandFromRequest(r *http.Request, parameter string) (*models.Command, error) {
	if r.URL.Query().Get(parameter) == "" {
		return nil, errors.New(parameter + " parameter not found")
	}

	id, err := strconv.Atoi(r.URL.Query().Get(parameter))
	if err != nil {
		return nil, err
	}

	command, err := app.commandRepository.Find(id)
	if err != nil {
		return nil, err
	}

	return command, nil
}

func (app *App) getCommandImages() []map[string]interface{} {
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
