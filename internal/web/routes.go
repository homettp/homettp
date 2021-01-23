package web

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *App) routes() http.Handler {
	baseMiddleware := alice.New(app.recoverPanic)
	webMiddleware := alice.New(
		app.sessionManager.LoadAndSave,
		app.remember,
		app.authenticate,
		app.flashMessage,
		app.inertiaManager.Middleware,
	)

	mux := http.NewServeMux()
	mux.Handle("/", webMiddleware.Append(app.redirectIfNotAuthenticated).ThenFunc(app.commandIndex))
	mux.Handle("/login", webMiddleware.Append(app.redirectIfAuthenticated).ThenFunc(app.login))
	mux.Handle("/logout", webMiddleware.Append(app.redirectIfNotAuthenticated).ThenFunc(app.logout))
	mux.Handle("/user", webMiddleware.Append(app.redirectIfNotAuthenticated).ThenFunc(app.userIndex))

	fileServer := http.FileServer(http.Dir("./public/"))

	mux.Handle("/css/", fileServer)
	mux.Handle("/images/", fileServer)
	mux.Handle("/js/", fileServer)
	mux.Handle("/favicon.ico", fileServer)

	return baseMiddleware.Then(mux)
}
