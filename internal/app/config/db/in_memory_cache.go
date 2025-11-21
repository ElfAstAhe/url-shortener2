package db

import (
	"sync"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/model"
)

type InMemoryCache interface {
	GetShortURIRWMutex() *sync.RWMutex
	GetShortURIUserRWMutex() *sync.RWMutex
	GetShortURIAuthRWMutex() *sync.RWMutex
	GetShortURICache() map[string]*model.ShortURI
	GetShortURIUserCache() map[string]*model.ShortURIUser
	GetShortURIAuditCache() map[string]*model.ShortURIAudit
}
