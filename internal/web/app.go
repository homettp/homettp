package web

import (
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/chmike/securecookie"
	"github.com/gomodule/redigo/redis"
	"github.com/homettp/homettp/internal/models"
	"github.com/petaki/inertia-go"
	"github.com/petaki/support-go/mix"
)

const (
	sessionKeyAuthUserId   = "authUserId"
	sessionKeyFlashMessage = "flashMessage"
	sessionKeyIntendedUrl  = "intendedUrl"
	rememberCookieName     = "remember"
)

type App struct {
	debug             bool
	url               string
	commandTimeout    int
	errorLog          *log.Logger
	infoLog           *log.Logger
	redisPool         *redis.Pool
	redisKeyPrefix    string
	sessionManager    *scs.SessionManager
	rememberCookie    *securecookie.Obj
	mixManager        *mix.Mix
	inertiaManager    *inertia.Inertia
	queue             <-chan int
	commandRepository models.CommandRepository
	callRepository    models.CallRepository
	userRepository    models.UserRepository
}
