package bootstrap

import (
	"github.com/ElfAstAhe/url-shortener2/internal/app/config/db"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
)

type App struct {
	db     db.DB
	logger logger.Logger
}
