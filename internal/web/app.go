package web

import (
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/chmike/securecookie"
	"github.com/homettp/homettp/internal/config"
	"github.com/homettp/homettp/internal/models"
	"github.com/petaki/inertia-go"
	"github.com/petaki/support-go/vite"
)

const (
	sessionKeyAuthUserID   = "authUserID"
	sessionKeyFlashMessage = "flashMessage"
	sessionKeyIntendedURL  = "intendedURL"
	rememberCookieName     = "remember"
)

type app struct {
	appConfig         *config.Config
	infoLog           *log.Logger
	errorLog          *log.Logger
	sessionManager    *scs.SessionManager
	rememberCookie    *securecookie.Obj
	viteManager       *vite.Vite
	inertiaManager    *inertia.Inertia
	queue             chan int64
	commandRepository models.CommandRepository
	callRepository    models.CallRepository
	userRepository    models.UserRepository
}
