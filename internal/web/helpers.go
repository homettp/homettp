package web

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/chmike/securecookie"
	"github.com/homettp/homettp/internal/models"
	"github.com/petaki/support-go/forms"
)

func (app *app) authUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(contextKeyAuthUser).(*models.User)
	if !ok {
		return nil
	}

	return user
}

func (app *app) generateToken() (string, error) {
	bytes, err := securecookie.GenerateRandomKey()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func (app *app) formError(w http.ResponseWriter, err error) {
	var fe *forms.Error

	if errors.As(err, &fe) {
		http.Error(w, fe.Msg, fe.Status)
	} else {
		app.serverError(w, err)
	}
}

func (app *app) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	if app.debug {
		http.Error(w, trace, http.StatusInternalServerError)
		return
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *app) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *app) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *app) methodNotAllowed(w http.ResponseWriter, allow []string) {
	w.Header().Set("Allow", strings.Join(allow, ", "))
	app.clientError(w, http.StatusMethodNotAllowed)
}
