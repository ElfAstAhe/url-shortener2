package db

import (
	"database/sql"
	"sync"

	"github.com/ElfAstAhe/url-shortener2/internal/app/config"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/model"
)

type inMemoryDB struct {
	DBKind string

	// ShortURI map key is short uri entity id attribute
	ShortURI map[string]*model.ShortURI

	// ShortURIAudit map key is short uri audit entity id attribute
	ShortURIAudit map[string]*model.ShortURIAudit

	// ShortURIUser map key is short uri user entity id attribute
	ShortURIUser map[string]*model.ShortURIUser

	shortURIAnchor      sync.RWMutex
	shortURIUserAnchor  sync.RWMutex
	shortURIAuditAnchor sync.RWMutex
}

var inMemDB *inMemoryDB

func newInMemoryDB() (*inMemoryDB, error) {
	if inMemDB != nil {
		return inMemDB, nil
	}

	inMemDB = &inMemoryDB{
		ShortURI:      make(map[string]*model.ShortURI),
		ShortURIAudit: make(map[string]*model.ShortURIAudit),
		ShortURIUser:  make(map[string]*model.ShortURIUser),
		DBKind:        config.DBKindInMemory,
	}

	return inMemDB, nil
}

// Closer

func (idb *inMemoryDB) Close() error {
	clear(idb.ShortURI)
	clear(idb.ShortURIAudit)
	clear(idb.ShortURIUser)

	return nil
}

// ========

// DB

func (idb *inMemoryDB) GetDB() *sql.DB {
	return nil
}

func (idb *inMemoryDB) GetDBKind() string {
	return idb.DBKind
}

func (idb *inMemoryDB) GetDsn() string {
	return ""
}

// ========

// InMemoryCache

func (idb *inMemoryDB) GetShortURIRWMutex() *sync.RWMutex {
	return &idb.shortURIAnchor
}

func (idb *inMemoryDB) GetShortURIUserRWMutex() *sync.RWMutex {
	return &idb.shortURIUserAnchor
}

func (idb *inMemoryDB) GetShortURIAuthRWMutex() *sync.RWMutex {
	return &idb.shortURIAuditAnchor
}

func (idb *inMemoryDB) GetShortURICache() map[string]*model.ShortURI {
	return idb.ShortURI
}

func (idb *inMemoryDB) GetShortURIAuditCache() map[string]*model.ShortURIAudit {
	return idb.ShortURIAudit
}

func (idb *inMemoryDB) GetShortURIUserCache() map[string]*model.ShortURIUser {
	return idb.ShortURIUser
}

// ========
