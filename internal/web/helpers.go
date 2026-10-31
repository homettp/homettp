package web

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"runtime/debug"
	"strconv"
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

func (a *app) isSameOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}

	originURL, err := url.Parse(origin)
	if err != nil {
		return false
	}

	if strings.EqualFold(originURL.Host, r.Host) {
		return true
	}

	appURL, err := url.Parse(a.appConfig.URL)

	return err == nil && appURL.Host != "" && strings.EqualFold(originURL.Host, appURL.Host)
}

func (a *app) generateToken() (string, error) {
	bytes, err := securecookie.GenerateRandomKey()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func (a *app) formError(w http.ResponseWriter, err error) {
	if fe, ok := errors.AsType[*forms.Error](err); ok {
		http.Error(w, fe.Msg, fe.Status)
	} else {
		a.serverError(w, err)
	}
}

func (a *app) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.errorLog.Output(2, trace)

	if a.appConfig.Debug {
		http.Error(w, trace, http.StatusInternalServerError)

		return
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *app) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (a *app) lookupError(w http.ResponseWriter, err error) {
	if errors.Is(err, models.ErrNoRecord) || errors.Is(err, ErrMissingParameter) || errors.Is(err, strconv.ErrSyntax) || errors.Is(err, strconv.ErrRange) {
		a.notFound(w)

		return
	}

	a.serverError(w, err)
}

func (a *app) notFound(w http.ResponseWriter) {
	a.clientError(w, http.StatusNotFound)
}

func (a *app) methodNotAllowed(w http.ResponseWriter, allow []string) {
	w.Header().Set("Allow", strings.Join(allow, ", "))
	a.clientError(w, http.StatusMethodNotAllowed)
}

func isSecureURL(rawURL string) bool {
	return strings.HasPrefix(strings.ToLower(rawURL), "https://")
}

func isKnownPath(r *http.Request) bool {
	return r.Pattern != "/" || r.URL.Path == "/"
}

func failedOutput(out []byte, err error) string {
	if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
		out = append(out, exitErr.Stderr...)
	}

	return string(out) + err.Error()
}
