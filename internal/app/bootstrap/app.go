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
	"github.com/ElfAstAhe/url-shortener2/internal/dal/repository"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/handler"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
	migrations "github.com/ElfAstAhe/url-shortener2/pkg/migrations/goose"
)

type App struct {
	ctx              context.Context
	cancelFunc       context.CancelFunc
	WG               sync.WaitGroup
	db               db.DB
	conf             *config.Config
	Log              logger.Logger
	shortURIUserRepo irepo.ShortURIUserRepository
	shortURIRepo     irepo.ShortURIRepository
	shorterService   service.Shorter
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
	log := app.Log.GetLogger("init")
	//    defer _utl.CloseOnly(logger.(io.Closer))

	log.Info("loading config")
	if err := app.loadConfig(); err != nil {
		return err
	}

	log.Info("initializing logger")
	if err := app.initLogger(); err != nil {
		return err
	}

	log.Info("initializing database")
	if err := app.initDatabase(); err != nil {
		return err
	}

	log.Info("migrate database")
	if err := app.migrateDatabase(); err != nil {
		return err
	}

	log.Info("initializing dependencies")
	if err := app.initDependencies(); err != nil {
		return err
	}

	log.Info("initializing startup services")
	if err := app.initStartupServices(); err != nil {
		return err
	}

	log.Info("initializing http server handlers")
	if err := app.initRouter(); err != nil {
		return err
	}

	log.Info("initializing http server")
	if err := app.initHTTPServer(); err != nil {
		return err
	}

	return nil
}

func (app *App) Run() error {
	log := app.Log.GetLogger("run")
	//    defer _utl.CloseOnly(logger.(io.Closer))
	log.Info("Starting graceful shutdown go routine...")
	app.WG.Add(1)
	go app.gracefulShutdown()

	log.Info("Starting server...")
	if err := app.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Errorf("Error starting server with error [%v]", err)

		return err
	}

	return nil
}

func (app *App) Close() error {
	if err := db.CloseDB(app.db); err != nil {
		return err
	}

	if err := app.Log.Close(); err != nil {
		return err
	}

	return nil
}

func (app *App) loadConfig() error {
	appConf := config.NewConfig()
	err := appConf.LoadConfig()
	if err != nil {
		return err
	}
	app.conf = appConf

	return nil
}

func (app *App) initLogger() error {
	fullLogger, err := logger.NewZapLogger(app.conf.LogLevel, "")
	if err != nil {
		return err
	}
	app.Log = fullLogger

	return nil
}

func (app *App) initDatabase() error {
	res, err := db.NewDB(app.db.GetDBKind(), app.conf.DBDsn)
	if err != nil {
		return err
	}
	app.db = res

	return nil
}

func (app *App) migrateDatabase() error {
	migrator, err := migrations.NewGooseDBMigrator(app.ctx, app.db.GetDB(), app.Log)
	if err != nil {
		return err
	}

	if err := migrator.Initialize(); err != nil {
		return err
	}

	if err := migrator.Up(); err != nil {
		return err
	}

	return nil
}

func (app *App) initDependencies() error {
	var err error
	// repositories
	if app.shortURIUserRepo, err = app.createShortURIUserRepo(); err != nil {
		return err
	}
	if app.shortURIRepo, err = app.createShortURIRepo(app.shortURIUserRepo); err != nil {
		return err
	}

	// services
	app.shorterService = service.NewShorterService(app.conf, app.shortURIRepo)

	// facades
	// ..
	// ..

	return nil
}

func (app *App) initStartupServices() error {
	// ..

	return nil
}

func (app *App) initRouter() error {
	// ToDo: implement
	//    app.router = handler.NewChiRouter(app.conf, .., app.log)

	return nil
}

func (app *App) initHTTPServer() error {
	app.httpServer = &http.Server{
		Addr:    app.conf.HTTP.GetListenerAddr(),
		Handler: app.router.GetRouter(),
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

func (app *App) createShortURIUserRepo() (irepo.ShortURIUserRepository, error) {
	if app.db.GetDBKind() == config.DBKindPostgres {
		return repository.NewShortURIUserPgRepo(app.db)
	}

	return repository.NewShortURIUserInMemRepo(app.db)
}

func (app *App) createShortURIRepo(shortURIUserRepo irepo.ShortURIUserRepository) (irepo.ShortURIRepository, error) {
	if app.db.GetDBKind() == config.DBKindPostgres {
		return repository.NewShortURIPgRepo(app.db, shortURIUserRepo)
	}

	return repository.NewShortURIInMemRepo(app.db, shortURIUserRepo)
}
