package pgdb

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//type Book struct {
//	Id     int    `json:"id" gorm:"primaryKey"`
//	Title  string `json:"title"`
//	Author string `json:"author"`
//	Desc   string `json:"desc"`
//}

type GoPg struct {
	DB     *gorm.DB
	URI    string
	DBName string
	TBName string
}

//type Querier interface {
//	FilterWithNameAndRole(firstName string) ([]gen.T, error)
//}

func New(uri, dbName, tbName, ssl string, mStruct interface{}) (repo *GoPg, err error) {
	dsn := fmt.Sprintf("%s/%s?sslmode=%s", uri, dbName, ssl)
	log.Println("dsn : ", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	repo = &GoPg{db, uri, dbName, tbName}

	//g := gen.NewGenerator(gen.Config{
	//	OutPath: "./query",
	//	Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	//})
	//
	//g.UseDB(db)
	//
	//g.ApplyBasic(staff.Staff{})
	//
	//g.Execute()

	err = repo.DB.AutoMigrate(&mStruct)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
