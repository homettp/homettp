package web

import "net/http"

func (app *App) userIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		app.methodNotAllowed(w, []string{"GET"})

		return
	}

	err := app.inertiaManager.Render(w, r, "user/Index", map[string]interface{}{
		"isUsersActive": true,
	})
	if err != nil {
		app.serverError(w, err)
	}
}
