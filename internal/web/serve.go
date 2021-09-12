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
	"github.com/homettp/homettp/internal/models"
	"github.com/homettp/homettp/resources/views"
	"github.com/petaki/inertia-go"
	"github.com/petaki/support-go/cli"
	"github.com/petaki/support-go/mix"
)

// Serve function.
func Serve(debug bool, addr, url, key, redisKeyPrefix string, redisPool *redis.Pool, commandTimeout, commandWorkerCount, commandHistoryLimit int) {
	sessionManager := scs.New()
	sessionManager.Store = redisstore.NewWithPrefix(redisPool, fmt.Sprintf("%sscs:session:", redisKeyPrefix))

	rememberCookie, err := securecookie.New(rememberCookieName, []byte(key), securecookie.Params{
		Path:     "/",
		MaxAge:   157680000, // Five years
		HTTPOnly: true,
	})
	if err != nil {
		cli.ErrorLog.Fatal(err)
	}

	mixManager, inertiaManager, err := newMixAndInertiaManager(url)
	if err != nil {
		cli.ErrorLog.Fatal(err)
	}

	queue := make(chan int64, 100)

	webApp := &app{
		debug:               debug,
		url:                 url,
		infoLog:             cli.InfoLog,
		errorLog:            cli.ErrorLog,
		redisPool:           redisPool,
		redisKeyPrefix:      redisKeyPrefix,
		commandTimeout:      commandTimeout,
		commandWorkerCount:  commandWorkerCount,
		commandHistoryLimit: commandHistoryLimit,
		sessionManager:      sessionManager,
		rememberCookie:      rememberCookie,
		mixManager:          mixManager,
		inertiaManager:      inertiaManager,
		queue:               queue,
		commandRepository: &models.RedisCommandRepository{
			RedisPool:      redisPool,
			RedisKeyPrefix: redisKeyPrefix,
		},
		callRepository: &models.RedisCallRepository{
			RedisPool:           redisPool,
			RedisKeyPrefix:      redisKeyPrefix,
			CommandHistoryLimit: commandHistoryLimit,
		},
		userRepository: &models.RedisUserRepository{
			RedisPool:      redisPool,
			RedisKeyPrefix: redisKeyPrefix,
		},
	}

	for i := 0; i < commandWorkerCount; i++ {
		go webApp.worker()
	}

	srv := &http.Server{
		Addr:         addr,
		ErrorLog:     webApp.errorLog,
		Handler:      webApp.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	webApp.infoLog.Printf("Starting server on "+cli.Green("%s"), addr)

	go func() {
		err = srv.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			webApp.errorLog.Fatal(err)
		}
	}()

	<-done
	webApp.infoLog.Print("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(webApp.commandTimeout)*time.Second)
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

func newMixAndInertiaManager(url string) (*mix.Mix, *inertia.Inertia, error) {
	mixManager := mix.New("")

	version, err := mixManager.Hash("")
	if err != nil {
		return nil, nil, err
	}

	inertiaManager := inertia.NewWithFS(url, "app.gohtml", version, views.Templates)
	inertiaManager.Share("title", "Homettp")
	inertiaManager.ShareFunc("mix", mixManager.Mix)

	return mixManager, inertiaManager, nil
}
