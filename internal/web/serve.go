package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/chmike/securecookie"
	"github.com/gomodule/redigo/redis"
	"github.com/homettp/homettp/internal/config"
	"github.com/homettp/homettp/internal/models"
	"github.com/homettp/homettp/resources/views"
	"github.com/homettp/homettp/static"
	"github.com/petaki/inertia-go"
	"github.com/petaki/support-go/cli"
	"github.com/petaki/support-go/vite"
	supportweb "github.com/petaki/support-go/web"
)

// Serve function.
func Serve(appConfig *config.Config, redisPool *redis.Pool) error {
	sessionManager := scs.New()
	sessionManager.Store = redisstore.NewWithPrefix(redisPool, fmt.Sprintf("%sscs:session:", appConfig.RedisKeyPrefix))
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = isSecureURL(appConfig.URL)

	rememberCookie, err := securecookie.New(rememberCookieName, []byte(appConfig.Key), securecookie.Params{
		Path:     "/",
		MaxAge:   157680000, // Five years
		HTTPOnly: true,
		Secure:   isSecureURL(appConfig.URL),
		SameSite: securecookie.Lax,
	})
	if err != nil {
		return err
	}

	inertiaManager, err := newInertiaManager(appConfig)
	if err != nil {
		return err
	}

	queue := make(chan int64, 100)

	webApp := &app{
		appConfig:      appConfig,
		infoLog:        cli.InfoLog,
		errorLog:       cli.ErrorLog,
		sessionManager: sessionManager,
		rememberCookie: rememberCookie,
		inertiaManager: inertiaManager,
		queue:          queue,
		commandRepository: &models.RedisCommandRepository{
			RedisPool:      redisPool,
			RedisKeyPrefix: appConfig.RedisKeyPrefix,
		},
		callRepository: &models.RedisCallRepository{
			RedisPool:      redisPool,
			RedisKeyPrefix: appConfig.RedisKeyPrefix,
			HistoryLimit:   appConfig.CommandHistoryLimit,
		},
		userRepository: &models.RedisUserRepository{
			RedisPool:      redisPool,
			RedisKeyPrefix: appConfig.RedisKeyPrefix,
		},
	}

	err = webApp.failInterruptedCalls()
	if err != nil {
		webApp.errorLog.Print(err)
	}

	workerCtx, cancelWorkers := context.WithCancel(context.Background())
	defer cancelWorkers()

	for range appConfig.CommandWorkerCount {
		webApp.workers.Go(func() {
			webApp.worker(workerCtx)
		})
	}

	srv := supportweb.NewServer(appConfig.Addr, webApp.routes(), webApp.errorLog)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err = srv.Run(ctx)
	stop()

	if err != nil {
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(appConfig.CommandTimeout)*time.Second)
	defer cancel()

	webApp.closeQueue()

	err = srv.Shutdown(shutdownCtx)
	if err != nil {
		webApp.errorLog.Print(err)
	}

	webApp.waitWorkers(shutdownCtx, cancelWorkers)

	webApp.infoLog.Print("Server exited properly")

	return nil
}

func newInertiaManager(appConfig *config.Config) (*inertia.Inertia, error) {
	var viteManager *vite.Vite
	var version string
	var err error

	if appConfig.Debug {
		viteManager = vite.New("static", "build")
	} else {
		viteManager = vite.New("static", "build", static.Files)
	}

	version, err = viteManager.ManifestHash()
	if err != nil {
		return nil, err
	}

	inertiaManager := inertia.New(appConfig.URL, "app.gohtml", version, views.Templates)
	inertiaManager.Share("title", "Homettp")

	for name, helper := range viteManager.Helpers() {
		inertiaManager.ShareFunc(name, helper)
	}

	return inertiaManager, nil
}
