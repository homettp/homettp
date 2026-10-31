package web

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/chmike/securecookie"
	"github.com/homettp/homettp/internal/config"
	"github.com/homettp/homettp/internal/models"
	"github.com/petaki/inertia-go"
)

type failingUserRepository struct {
	models.UserRepository
}

func (failingUserRepository) Find(int) (*models.User, error) {
	return nil, errors.New("redis: i/o timeout")
}

func TestPreventCSRF(t *testing.T) {
	a := &app{appConfig: &config.Config{URL: "https://home.example.com"}}
	h := a.preventCSRF(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	tests := []struct {
		name    string
		method  string
		inertia bool
		origin  string
		status  int
	}{
		{"get", "GET", false, "https://evil.example.com", http.StatusOK},
		{"post without inertia", "POST", false, "", http.StatusForbidden},
		{"delete cross origin", "DELETE", true, "https://evil.example.com", http.StatusForbidden},
		{"put null origin", "PUT", true, "null", http.StatusForbidden},
		{"post same host", "POST", true, "http://example.com", http.StatusOK},
		{"post app url", "POST", true, "https://home.example.com", http.StatusOK},
		{"post without origin", "POST", true, "", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, "http://example.com/command/delete?id=1", nil)
			if tt.inertia {
				r.Header.Set(inertia.HeaderInertia, "true")
			}

			if tt.origin != "" {
				r.Header.Set("Origin", tt.origin)
			}

			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)

			if w.Code != tt.status {
				t.Fatalf("expected %d, got %d", tt.status, w.Code)
			}
		})
	}
}

func TestAuthenticateRepositoryError(t *testing.T) {
	sessionManager := scs.New()
	a := &app{
		appConfig:      &config.Config{},
		errorLog:       log.New(io.Discard, "", 0),
		sessionManager: sessionManager,
		userRepository: failingUserRepository{},
	}

	setup := sessionManager.LoadAndSave(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionManager.Put(r.Context(), sessionKeyAuthUserID, 1)
	}))

	w := httptest.NewRecorder()
	setup.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))

	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(w.Result().Cookies()[0])

	w = httptest.NewRecorder()
	sessionManager.LoadAndSave(a.authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	}))).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestRememberInvalidCookie(t *testing.T) {
	rememberCookie, err := securecookie.New(rememberCookieName, []byte("0123456789abcdef0123456789abcdef"), securecookie.Params{Path: "/"})
	if err != nil {
		t.Fatal(err)
	}

	sessionManager := scs.New()
	a := &app{
		appConfig:      &config.Config{},
		errorLog:       log.New(io.Discard, "", 0),
		sessionManager: sessionManager,
		rememberCookie: rememberCookie,
	}

	for _, value := range []string{"", "garbage", "MTIz"} {
		t.Run(value, func(t *testing.T) {
			called := false
			h := sessionManager.LoadAndSave(a.remember(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
			})))

			r := httptest.NewRequest("GET", "/login", nil)
			r.AddCookie(&http.Cookie{Name: rememberCookieName, Value: value})

			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)

			if w.Code != http.StatusOK || !called {
				t.Fatalf("expected the request to continue as a guest, got %d", w.Code)
			}

			if !strings.Contains(w.Header().Get("Set-Cookie"), rememberCookieName+"=; Path=/; Max-Age=0") {
				t.Fatalf("expected the remember cookie to be deleted, got %q", w.Header().Values("Set-Cookie"))
			}
		})
	}
}

func TestIsKnownPath(t *testing.T) {
	var known bool

	mux := http.NewServeMux()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		known = isKnownPath(r)
	})

	mux.Handle("/", handler)
	mux.Handle("/user", handler)

	for path, expected := range map[string]bool{
		"/":            true,
		"/user":        true,
		"/.well-known": false,
		"/missing":     false,
	} {
		mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, path, nil))

		if known != expected {
			t.Errorf("%s: expected %v, got %v", path, expected, known)
		}
	}
}

func TestLookupError(t *testing.T) {
	a := &app{
		appConfig: &config.Config{},
		errorLog:  log.New(io.Discard, "", 0),
	}

	_, atoiErr := strconv.Atoi("x")

	for err, expected := range map[error]int{
		models.ErrNoRecord:               http.StatusNotFound,
		ErrMissingParameter:              http.StatusNotFound,
		atoiErr:                          http.StatusNotFound,
		errors.New("redis: i/o timeout"): http.StatusInternalServerError,
	} {
		w := httptest.NewRecorder()
		a.lookupError(w, err)

		if w.Code != expected {
			t.Errorf("%v: expected %d, got %d", err, expected, w.Code)
		}
	}
}
