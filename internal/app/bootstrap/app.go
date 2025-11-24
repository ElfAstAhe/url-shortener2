package bootstrap

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/ElfAstAhe/url-shortener2/internal/app/config"
	"github.com/ElfAstAhe/url-shortener2/internal/app/config/db"
	irepo "github.com/ElfAstAhe/url-shortener2/internal/bll/repository"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/service"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/audit"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/facade"
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
	connCheckRepo    irepo.DBConnCheckRepository
	shortURIUserRepo irepo.ShortURIUserRepository
	shortURIRepo     irepo.ShortURIRepository
	shorterService   service.Shorter
	toolFacade       facade.ToolFacade
	authFacade       facade.AuthFacade
	shortenFacade    facade.ShortenFacade
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
	if app.connCheckRepo, err = app.createConnCheckRepo(); err != nil {
		return err
	}
	if app.shortURIUserRepo, err = app.createShortURIUserRepo(); err != nil {
		return err
	}
	if app.shortURIRepo, err = app.createShortURIRepo(app.shortURIUserRepo); err != nil {
		return err
	}

	// services
	app.shorterService = service.NewShorterService(app.conf, app.shortURIRepo)

	// facades
	app.toolFacade = facade.NewToolFacadeImpl(app.connCheckRepo)
	//	app.authFacade = facade.NewAuthFacadeImpl()
	//	app.shortenFacade = facade.NewShortenFacadeImpl()

	return nil
}

func (app *App) initStartupServices() error {
	// ..

	return nil
}

func (app *App) initRouter() error {
	// observers
	observers, err := app.initIncomeObservers()
	if err != nil {
		return err
	}

	app.router = handler.NewAppChiRouter(app.toolFacade, app.authFacade, app.shortenFacade, observers, app.conf, app.Log)

	return nil
}

func (app *App) initIncomeObservers() ([]audit.IncomeObserver, error) {
	res := make([]audit.IncomeObserver, 0)
	// local
	if strings.TrimSpace(app.conf.AuditFile) != "" {
		localIncomeObserver, err := audit.NewIncomeLocalService(app.conf)
		if err != nil {
			return nil, err
		}

		res = append(res, localIncomeObserver)
	}
	// remote
	if strings.TrimSpace(app.conf.AuditURL) != "" {
		res = append(res, audit.NewIncomeRemoteService(app.conf))
	}

	return res, nil
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
