package bootstrap

import (
	"github.com/ElfAstAhe/url-shortener2/internal/app/config"
	"github.com/ElfAstAhe/url-shortener2/internal/app/config/db"
	irepo "github.com/ElfAstAhe/url-shortener2/internal/bll/repository"
	"github.com/ElfAstAhe/url-shortener2/internal/dal/repository"
	"github.com/ElfAstAhe/url-shortener2/internal/dal/storage"
	"github.com/ElfAstAhe/url-shortener2/internal/utils"
)

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
