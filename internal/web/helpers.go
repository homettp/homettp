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

func (a *app) authUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(contextKeyAuthUser).(*models.User)
	if !ok {
		return nil
	}

	return user
}

func (a *app) generateToken() (string, error) {
	bytes, err := securecookie.GenerateRandomKey()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func (a *app) formError(w http.ResponseWriter, err error) {
	var fe *forms.Error

	if errors.As(err, &fe) {
		http.Error(w, fe.Msg, fe.Status)
	} else {
		a.serverError(w, err)
	}
}

func (a *app) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.errorLog.Output(2, trace)

	if a.debug {
		http.Error(w, trace, http.StatusInternalServerError)
		return
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *app) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (a *app) notFound(w http.ResponseWriter) {
	a.clientError(w, http.StatusNotFound)
}

func (a *app) methodNotAllowed(w http.ResponseWriter, allow []string) {
	w.Header().Set("Allow", strings.Join(allow, ", "))
	a.clientError(w, http.StatusMethodNotAllowed)
}
