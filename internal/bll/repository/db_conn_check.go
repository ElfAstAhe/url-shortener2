package repository

type DBConnCheckRepository interface {
	CheckDBConn() error
	Close() error
}
