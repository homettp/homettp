package web

import (
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/chmike/securecookie"
	"github.com/homettp/homettp/internal/models"
	"github.com/petaki/inertia-go"
	"github.com/petaki/support-go/mix"
)

const (
	sessionKeyAuthUserID   = "authUserID"
	sessionKeyFlashMessage = "flashMessage"
	sessionKeyIntendedURL  = "intendedURL"
	rememberCookieName     = "remember"
)

type app struct {
	debug             bool
	url               string
	infoLog           *log.Logger
	errorLog          *log.Logger
	commandTimeout    int
	sessionManager    *scs.SessionManager
	rememberCookie    *securecookie.Obj
	mixManager        *mix.Mix
	inertiaManager    *inertia.Inertia
	queue             chan int64
	commandRepository models.CommandRepository
	callRepository    models.CallRepository
	userRepository    models.UserRepository
}
