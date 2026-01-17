package bootstrap

import (
	"net/http"
	"strings"

	"github.com/ElfAstAhe/url-shortener2/internal/app/config"
	"github.com/ElfAstAhe/url-shortener2/internal/app/config/db"
	irepo "github.com/ElfAstAhe/url-shortener2/internal/bll/repository"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/service"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/audit"
	"github.com/ElfAstAhe/url-shortener2/internal/dal/repository"
	"github.com/ElfAstAhe/url-shortener2/internal/dal/storage"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/facade"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/handler"
	"github.com/ElfAstAhe/url-shortener2/internal/utils"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
	migrations "github.com/ElfAstAhe/url-shortener2/pkg/migrations/goose"
)

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
	res, err := db.NewDB(app.conf.DBKind, app.conf.DBDsn)
	if err != nil {
		return err
	}
	app.db = res

	return nil
}

func (app *App) migrateDatabase() error {
	if app.db.GetDBKind() == config.DBKindPostgres {
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
	}

	return nil
}

func (app *App) loadInMemData() error {
	if cache, ok := app.db.(db.InMemoryCache); ok {
		log := app.Log.GetLogger("inMem")
		log.Info("Load data from storage...")
		if err := app.loadShortURIData(app.conf.StoragePath, cache); err != nil {
			log.Errorf("Error loading data: [%v]", err)
			log.Warn("Using empty data storage")
		}

		log.Info("Load data from storage user...")
		if err := app.loadShortURIUserData(app.conf.StorageUserPath, cache); err != nil {
			log.Errorf("Error loading data: [%v]", err)
			log.Warn("Using empty data storage")
		}
	}

	return nil
}

func (app *App) saveInMemData() error {
	log := app.Log.GetLogger("inMem")
	if cache, ok := app.db.(db.InMemoryCache); ok {
		if err := app.saveShortURIData(config.AppConfig.StoragePath, cache); err != nil {
			log.Errorf("error save shortURI data: [%v]", err)
		}
		if err := app.saveShortURIUserData(config.AppConfig.StorageUserPath, cache); err != nil {
			log.Errorf("error save shortURI user data: [%v]", err)
		}
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
	app.shortenFacade = facade.NewShortenFacadeImpl(app.shorterService)
	app.userFacade = facade.NewUserFacadeImpl(app.shorterService, app.Log)

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

	app.router = handler.NewAppChiRouter(app.toolFacade, app.shortenFacade, app.userFacade, observers, app.conf, app.Log)

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

func (app *App) createConnCheckRepo() (irepo.DBConnCheckRepository, error) {
	if app.db.GetDBKind() == config.DBKindPostgres {
		return repository.NewDBConnCheckPgRepo(app.db)
	}

	return repository.NewDBConnCheckImMemRepo()
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

func (app *App) loadShortURIData(storagePath string, cache db.InMemoryCache) error {
	storageReader, err := storage.NewShortURLStorageReader(storagePath, app.Log)
	if err != nil {
		return err
	}
	defer utils.CloseOnly(storageReader)

	return storageReader.LoadData(cache.GetShortURICache())
}

func (app *App) loadShortURIUserData(storagePath string, cache db.InMemoryCache) error {
	storageReader, err := storage.NewShortURIUserStorageReader(storagePath, app.Log)
	if err != nil {
		return err
	}
	defer utils.CloseOnly(storageReader)

	return storageReader.LoadData(cache.GetShortURIUserCache())
}

func (app *App) saveShortURIData(storagePath string, cache db.InMemoryCache) error {
	storageWriter, err := storage.NewShortURLStorageWriter(storagePath, app.Log)
	if err != nil {
		return err
	}
	defer utils.CloseOnly(storageWriter)

	return storageWriter.SaveData(cache.GetShortURICache())
}

func (app *App) saveShortURIUserData(storagePath string, cache db.InMemoryCache) error {
	storageWriter, err := storage.NewShortURIUserStorageWriter(storagePath, app.Log)
	if err != nil {
		return err
	}
	defer utils.CloseOnly(storageWriter)

	return storageWriter.SaveData(cache.GetShortURIUserCache())
}
