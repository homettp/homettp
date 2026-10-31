package web

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
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
		a.lookupError(w, err)

		return
	}

	if subtle.ConstantTimeCompare([]byte(command.Token), []byte(r.URL.Query().Get("token"))) != 1 {
		a.notFound(w)

		return
	}

	payload := r.URL.Query().Get("payload")

	if !isSafePayload(payload) {
		a.clientError(w, http.StatusUnprocessableEntity)

		return
	}

	call := models.NewCall(command)
	call.Payload = payload

	err = a.callRepository.Create(call)
	if errors.Is(err, models.ErrNoRecord) {
		a.notFound(w)

		return
	}

	if err != nil {
		a.serverError(w, err)

		return
	}

	err = a.enqueue(call.ID)
	if err != nil {
		err = a.callRepository.Update(call, &models.Call{
			Status: models.Failed,
			Output: err.Error(),
		})
		if err != nil {
			a.serverError(w, err)

			return
		}

		a.clientError(w, http.StatusServiceUnavailable)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(call)
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *app) callHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		a.methodNotAllowed(w, []string{"GET", "HEAD"})

		return
	}

	command, err := a.commandFromRequest(r, "commandID")
	if err != nil {
		a.lookupError(w, err)

		return
	}

	calls, err := a.callRepository.FindAllByCommand(command)
	if err != nil {
		a.serverError(w, err)

		return
	}

	err = a.inertiaManager.Render(w, r, "call/History", map[string]any{
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
		a.lookupError(w, err)

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
		return nil, fmt.Errorf("%w: %s", ErrMissingParameter, parameter)
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
