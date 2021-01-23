package web

import "net/http"

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
