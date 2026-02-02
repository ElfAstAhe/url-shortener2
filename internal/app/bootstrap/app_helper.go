package bootstrap

import (
	"net/http"

	pb "github.com/ElfAstAhe/url-shortener2/api/proto"
	"github.com/ElfAstAhe/url-shortener2/internal/app/config"
	"github.com/ElfAstAhe/url-shortener2/internal/app/config/db"
	irepo "github.com/ElfAstAhe/url-shortener2/internal/bll/repository"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/service"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/audit"
	"github.com/ElfAstAhe/url-shortener2/internal/dal/repository"
	"github.com/ElfAstAhe/url-shortener2/internal/dal/storage"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/facade"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/handler"
	appgrpc "github.com/ElfAstAhe/url-shortener2/internal/grpc"
	"github.com/ElfAstAhe/url-shortener2/internal/utils"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
	migrations "github.com/ElfAstAhe/url-shortener2/pkg/migrations/goose"
	"google.golang.org/grpc"
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
		log.Info("Save data to storage...")
		if err := app.saveShortURIData(app.conf.StoragePath, cache); err != nil {
			log.Errorf("error save shortURI data: [%v]", err)
		}
		log.Info("Save data to storage user...")
		if err := app.saveShortURIUserData(app.conf.StorageUserPath, cache); err != nil {
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
	app.auditEventService = audit.NewIncomeEventService(app.Log)
	observers, err := app.initIncomeObservers()
	if err != nil {
		return err
	}
	for _, observer := range observers {
		app.auditEventService.Register(observer)
	}

	// facades
	app.toolFacade = facade.NewToolFacadeImpl(app.connCheckRepo, app.shorterService, app.conf.TrustedSubnetCIDR)
	app.shortenFacade = facade.NewShortenFacadeImpl(app.shorterService, app.conf.BaseURL)
	app.userFacade = facade.NewUserFacadeImpl(app.shorterService, app.Log)

	// grpc
	// facade
	app.grpcShortenFacade = appgrpc.NewShortenGRPCFacade(app.shorterService, app.conf.BaseURL, app.Log)

	// service
	app.grpcService = appgrpc.NewAppGRPCService(app.grpcShortenFacade, app.conf, app.Log)

	return nil
}

func (app *App) initStartupServices() error {
	// nothing

	return nil
}

func (app *App) initRouter() error {
	app.router = handler.NewAppChiRouter(app.toolFacade, app.shortenFacade, app.userFacade, app.auditEventService, app.conf, app.Log)

	return nil
}

func (app *App) initIncomeObservers() ([]audit.IncomeObserver, error) {
	// local
	localIncomeObserver, err := audit.NewIncomeLocalService(app.conf.AuditFile, app.conf.AuditIncomeLocal)
	if err != nil {
		return nil, err
	}
	// remote
	remoteIncomeObserver := audit.NewIncomeRemoteService(app.conf.AuditURL, app.conf.AuditIncomeRemote)

	return []audit.IncomeObserver{
		localIncomeObserver,
		remoteIncomeObserver,
	}, nil
}

func (app *App) initHTTPServer() error {
	app.httpServer = &http.Server{
		Addr:    app.conf.HTTP.GetListenerAddr(),
		Handler: app.router.GetRouter(),
	}

	return nil
}

func (app *App) initGRPCServer() error {
	app.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(appgrpc.NewAuthIter14Interceptor([]string{
			"ShortenURL",
		}, app.Log).UnaryInterceptor),
		grpc.UnaryInterceptor(appgrpc.NewAuthRetrieveInterceptor(app.Log).UnaryInterceptor),
		grpc.UnaryInterceptor(appgrpc.NewAuthTrailerInterceptor(app.Log).UnaryInterceptor),
	)
	pb.RegisterShortenerServiceServer(app.grpcServer, app.grpcService)

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
