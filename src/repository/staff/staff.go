package staff

import (
	"github.com/wayla99/go_gorm.git/src/repository/pgdb"
)

type Repository struct {
	*pgdb.GoPg
}

func New(uri, dbName, tbName, ssl string, mStruct interface{}) (repo *Repository, err error) {
	db, err := pgdb.New(uri, dbName, tbName, ssl, mStruct)
	if err != nil {
		return nil, err
	}
	return &Repository{db}, nil
}
