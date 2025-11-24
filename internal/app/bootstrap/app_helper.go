package bootstrap

import (
	"github.com/ElfAstAhe/url-shortener2/internal/app/config"
	irepo "github.com/ElfAstAhe/url-shortener2/internal/bll/repository"
	"github.com/ElfAstAhe/url-shortener2/internal/dal/repository"
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
