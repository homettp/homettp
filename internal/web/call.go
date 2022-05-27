package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/homettp/homettp/internal/models"
)

func (a *app) callIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "POST" {
		a.methodNotAllowed(w, []string{"GET", "POST"})

		return
	}

	command, err := a.commandFromRequest(r, "commandID")
	if err != nil {
		a.notFound(w)

		return
	}

	if command.Token != r.URL.Query().Get("token") {
		a.notFound(w)

		return
	}

	call := models.NewCall(command)
	call.Payload = r.URL.Query().Get("payload")

	err = a.callRepository.Create(call)
	if err != nil {
		a.serverError(w, err)

		return
	}

	select {
	case a.queue <- call.ID:
		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(call)
		if err != nil {
			a.serverError(w, err)
		}
	default:
		a.clientError(w, 503)
	}
}

func (a *app) callHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		a.methodNotAllowed(w, []string{"GET"})

		return
	}

	command, err := a.commandFromRequest(r, "id")
	if err != nil {
		a.notFound(w)

		return
	}

	calls, err := a.callRepository.FindAllByCommand(command)
	if err != nil {
		a.serverError(w, err)

		return
	}

	err = a.inertiaManager.Render(w, r, "call/History", map[string]interface{}{
		"command": command,
		"calls":   calls,
	})
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) callDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		a.methodNotAllowed(w, []string{"DELETE"})

		return
	}

	call, err := a.callFromRequest(r, "id")
	if err != nil {
		a.notFound(w)

		return
	}

	err = a.callRepository.Delete(call)
	if err != nil {
		a.serverError(w, err)

		return
	}

	a.sessionManager.Put(r.Context(), sessionKeyFlashMessage, "Deleted successfully.")
	http.Redirect(w, r, fmt.Sprintf("/call/history?commandID=%v", call.CommandID), http.StatusSeeOther)
}

func (a *app) callFromRequest(r *http.Request, parameter string) (*models.Call, error) {
	if r.URL.Query().Get(parameter) == "" {
		return nil, fmt.Errorf("%s parameter not found", parameter)
	}

	id, err := strconv.ParseInt(r.URL.Query().Get(parameter), 10, 64)
	if err != nil {
		return nil, err
	}

	call, err := a.callRepository.Find(id)
	if err != nil {
		return nil, err
	}

	return call, nil
}
