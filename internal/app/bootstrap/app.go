package bootstrap

import (
	"context"

	"github.com/ElfAstAhe/url-shortener2/internal/app/config"
	"github.com/ElfAstAhe/url-shortener2/internal/app/config/db"
	irepo "github.com/ElfAstAhe/url-shortener2/internal/bll/repository"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
)

type App struct {
	ctx             context.Context
	db              db.DB
	logger          logger.Logger
	shorURIUserRepo irepo.ShortURIUserRepository
	shortURIRepo    irepo.ShortURIRepository
}

func NewApp(config config.Config) *App {
	return &App{}
}

func (app *App) Initialize() error {
	// ToDo: implement

	return nil
}

func (app *App) Run() error {
	// ToDo: implement

	return nil
}
