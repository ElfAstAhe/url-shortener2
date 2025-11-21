package repository

type DBConnCheckImMemRepo struct {
}

func NewDBConnCheckImMemRepo() (*DBConnCheckImMemRepo, error) {
	return &DBConnCheckImMemRepo{}, nil
}

func (D *DBConnCheckImMemRepo) CheckDBConn() error {
	//    return errs.NewDalDBConnCheckError("IN_MEMORY DB, where is no connection")
	return nil
}

func (D *DBConnCheckImMemRepo) Close() error {
	return nil
}
