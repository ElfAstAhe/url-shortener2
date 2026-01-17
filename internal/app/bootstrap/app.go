package bootstrap

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ElfAstAhe/url-shortener2/internal/app/config"
	"github.com/ElfAstAhe/url-shortener2/internal/app/config/db"
	irepo "github.com/ElfAstAhe/url-shortener2/internal/bll/repository"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/service"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/facade"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/handler"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
)

type App struct {
	ctx              context.Context
	cancelFunc       context.CancelFunc
	WG               sync.WaitGroup
	db               db.DB
	conf             *config.Config
	Log              logger.Logger
	connCheckRepo    irepo.DBConnCheckRepository
	shortURIUserRepo irepo.ShortURIUserRepository
	shortURIRepo     irepo.ShortURIRepository
	shorterService   service.Shorter
	toolFacade       facade.ToolFacade
	shortenFacade    facade.ShortenFacade
	userFacade       facade.UserFacade
	router           handler.AppRouter
	httpServer       *http.Server
}

func NewApp() *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		ctx:        ctx,
		cancelFunc: cancel,
		Log:        logger.NewStartupZapLogger(),
	}
}

func (app *App) Init() error {
	log := app.Log.GetLogger("bootstrap init")
	//    defer _utl.CloseOnly(logger.(io.Closer))

	log.Info("loading config")
	if err := app.loadConfig(); err != nil {
		return err
	}

	log.Info("init logger")
	if err := app.initLogger(); err != nil {
		return err
	}

	log.Info("init database")
	if err := app.initDatabase(); err != nil {
		return err
	}

	log.Info("migrate database")
	if err := app.migrateDatabase(); err != nil {
		return err
	}

	log.Info("load im mem data")
	if err := app.loadInMemData(); err != nil {
		return err
	}

	log.Info("init dependencies")
	if err := app.initDependencies(); err != nil {
		return err
	}

	log.Info("init startup services")
	if err := app.initStartupServices(); err != nil {
		return err
	}

	log.Info("init http router")
	if err := app.initRouter(); err != nil {
		return err
	}

	log.Info("init http server")
	if err := app.initHTTPServer(); err != nil {
		return err
	}

	return nil
}

func (app *App) Run() error {
	log := app.Log.GetLogger("bootstrap run")
	//    defer _utl.CloseOnly(logger.(io.Closer))
	log.Info("start graceful shutdown go routine...")
	app.WG.Add(1)
	go app.gracefulShutdown()

	log.Info("start server...")
	if err := app.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Errorf("Error starting server with error [%v]", err)

		return err
	}

	return nil
}

func (app *App) Close() error {
	log := app.Log.GetLogger("bootstrap close")

	log.Info("close db connection")
	if err := db.CloseDB(app.db); err != nil {
		return err
	}

	log.Info("save in mem data")
	if err := app.saveInMemData(); err != nil {
		return err
	}

	return nil
}

func (app *App) gracefulShutdown() {
	// channel
	sig := make(chan os.Signal, 1)
	// register channel signals
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	// awaiting signal
	select {
	case <-sig:
		{
			app.cancelFunc()
			break
		}
	case <-app.ctx.Done():
		{
			signal.Stop(sig)
			break
		}
	}

	if err := app.httpServer.Shutdown(context.Background()); err != nil {
		app.Log.Errorf("error graceful shutdown http server with error [%v]", err)
	}

	app.WG.Done()
}
