package web

import "net/http"

func (app *App) commandIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		app.methodNotAllowed(w, []string{"GET"})

		return
	}

	err := app.inertiaManager.Render(w, r, "command/Index", map[string]interface{}{
		"isCommandsActive": true,
	})
	if err != nil {
		app.serverError(w, err)
	}
}
