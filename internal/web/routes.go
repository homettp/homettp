package web

import (
	"net/http"

	"github.com/justinas/alice"
)

func (a *app) routes() http.Handler {
	baseMiddleware := alice.New(a.recoverPanic)
	webMiddleware := alice.New(
		a.sessionManager.LoadAndSave,
		a.remember,
		a.authenticate,
		a.flashMessage,
		a.inertiaManager.Middleware,
	)

	mux := http.NewServeMux()
	mux.Handle("/", webMiddleware.Append(a.redirectIfNotAuthenticated).ThenFunc(a.commandIndex))
	mux.Handle("/command/create", webMiddleware.Append(a.redirectIfNotAuthenticated).ThenFunc(a.commandCreate))
	mux.Handle("/command/edit", webMiddleware.Append(a.redirectIfNotAuthenticated).ThenFunc(a.commandEdit))
	mux.Handle("/command/refresh-token", webMiddleware.Append(a.redirectIfNotAuthenticated).ThenFunc(a.commandRefreshToken))
	mux.Handle("/command/delete", webMiddleware.Append(a.redirectIfNotAuthenticated).ThenFunc(a.commandDelete))
	mux.Handle("/call", baseMiddleware.ThenFunc(a.callIndex))
	mux.Handle("/call/history", webMiddleware.Append(a.redirectIfNotAuthenticated).ThenFunc(a.callHistory))
	mux.Handle("/call/delete", webMiddleware.Append(a.redirectIfNotAuthenticated).ThenFunc(a.callDelete))
	mux.Handle("/user/create", webMiddleware.Append(a.redirectIfNotAuthenticated).ThenFunc(a.userCreate))
	mux.Handle("/user/edit", webMiddleware.Append(a.redirectIfNotAuthenticated).ThenFunc(a.userEdit))
	mux.Handle("/user/delete", webMiddleware.Append(a.redirectIfNotAuthenticated).ThenFunc(a.userDelete))
	mux.Handle("/login", webMiddleware.Append(a.redirectIfAuthenticated).ThenFunc(a.login))
	mux.Handle("/logout", webMiddleware.Append(a.redirectIfNotAuthenticated).ThenFunc(a.logout))
	mux.Handle("/user", webMiddleware.Append(a.redirectIfNotAuthenticated).ThenFunc(a.userIndex))

	fileServer := http.FileServer(http.Dir("./public/"))

	mux.Handle("/css/", fileServer)
	mux.Handle("/images/", fileServer)
	mux.Handle("/js/", fileServer)
	mux.Handle("/favicon.ico", fileServer)

	return baseMiddleware.Then(mux)
}
