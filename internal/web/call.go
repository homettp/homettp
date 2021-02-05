package web

import (
	"encoding/json"
	"net/http"

	"github.com/homettp/homettp/internal/models"
)

func (app *App) callIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "POST" {
		app.methodNotAllowed(w, []string{"GET", "POST"})

		return
	}

	command, err := app.getCommandFromRequest(r, "id")
	if err != nil {
		app.notFound(w)

		return
	}

	if command.Token != r.URL.Query().Get("token") {
		app.notFound(w)

		return
	}

	call := models.NewCall(command)
	call.Payload = r.URL.Query().Get("payload")

	err = app.callRepository.Create(call)
	if err != nil {
		app.serverError(w, err)

		return
	}

	select {
	case app.queue <- call.Id:
		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(call)
		if err != nil {
			app.serverError(w, err)
		}
	default:
		app.clientError(w, 503)
	}
}

func (app *App) callHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		app.methodNotAllowed(w, []string{"GET"})

		return
	}

	command, err := app.getCommandFromRequest(r, "id")
	if err != nil {
		app.notFound(w)

		return
	}

	calls, err := app.callRepository.FindAllByCommand(command)
	if err != nil {
		app.serverError(w, err)

		return
	}

	err = app.inertiaManager.Render(w, r, "call/History", map[string]interface{}{
		"command": command,
		"calls":   calls,
	})
	if err != nil {
		app.serverError(w, err)
	}
}
