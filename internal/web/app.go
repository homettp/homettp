package web

import (
	"log"
	"sync"

	"github.com/alexedwards/scs/v2"
	"github.com/chmike/securecookie"
	"github.com/homettp/homettp/internal/config"
	"github.com/homettp/homettp/internal/models"
	"github.com/petaki/inertia-go"
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
	inertiaManager    *inertia.Inertia
	queue             chan int64
	queueMu           sync.RWMutex
	queueClosed       bool
	workers           sync.WaitGroup
	commandRepository models.CommandRepository
	callRepository    models.CallRepository
	userRepository    models.UserRepository
}
