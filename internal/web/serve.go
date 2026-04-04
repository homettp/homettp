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
)

// Serve function.
func Serve(appConfig *config.Config, redisPool *redis.Pool) {
	sessionManager := scs.New()
	sessionManager.Store = redisstore.NewWithPrefix(redisPool, fmt.Sprintf("%sscs:session:", appConfig.RedisKeyPrefix))

	rememberCookie, err := securecookie.New(rememberCookieName, []byte(appConfig.Key), securecookie.Params{
		Path:     "/",
		MaxAge:   157680000, // Five years
		HTTPOnly: true,
	})
	if err != nil {
		cli.ErrorLog.Fatal(err)
	}

	viteManager, inertiaManager, err := newViteAndInertiaManager(appConfig)
	if err != nil {
		cli.ErrorLog.Fatal(err)
	}

	queue := make(chan int64, 100)

	webApp := &app{
		appConfig:      appConfig,
		infoLog:        cli.InfoLog,
		errorLog:       cli.ErrorLog,
		sessionManager: sessionManager,
		rememberCookie: rememberCookie,
		viteManager:    viteManager,
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

	for range appConfig.CommandWorkerCount {
		go webApp.worker()
	}

	srv := &http.Server{
		Addr:         appConfig.Addr,
		ErrorLog:     webApp.errorLog,
		Handler:      webApp.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	webApp.infoLog.Printf("Starting server on "+cli.Green("%s"), appConfig.Addr)

	go func() {
		err = srv.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			webApp.errorLog.Fatal(err)
		}
	}()

	<-done
	webApp.infoLog.Print("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(appConfig.CommandTimeout)*time.Second)
	defer func() {
		close(queue)
		cancel()
	}()

	err = srv.Shutdown(ctx)
	if err != nil {
		webApp.errorLog.Fatal(err)
	}

	webApp.infoLog.Print("Server exited properly")
}

func newViteAndInertiaManager(appConfig *config.Config) (*vite.Vite, *inertia.Inertia, error) {
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
		return nil, nil, err
	}

	inertiaManager := inertia.New(appConfig.URL, "app.gohtml", version, views.Templates)
	inertiaManager.Share("title", "Homettp")
	inertiaManager.ShareFunc("isRunningHot", viteManager.IsRunningHot)
	inertiaManager.ShareFunc("asset", viteManager.Asset)
	inertiaManager.ShareFunc("css", viteManager.CSS)

	return viteManager, inertiaManager, nil
}
