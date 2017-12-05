package persistence

import (
	"fmt"
	"os"
	"strconv"
	"database/sql"
	"github.com/ld100/goblet/util/log"
)

type DataSource struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func (ds *DataSource) DSN() string {
	var dsn string
	if len(ds.Database) > 0 {
		dsn = fmt.Sprintf(
			"host=%v user=%v dbname=%v sslmode=disable password=%v port=%v",
			ds.Host,
			ds.Username,
			ds.Database,
			ds.Password,
			ds.Port,
		)
	} else {
		dsn = ds.ShortDSN()
	}
	return dsn
}

func (ds *DataSource) ShortDSN() string {
	var dsn string
	dsn = fmt.Sprintf(
		"host=%v user=%v sslmode=disable password=%v port=%v",
		ds.Host,
		ds.Username,
		ds.Password,
		ds.Port,
	)
	return dsn
}

// Fetch data from environment variables
func (ds *DataSource) FetchENV() {
	ds.Host = os.Getenv("DB_HOST")
	ds.Port,_ = strconv.Atoi(os.Getenv("DB_PORT"))
	ds.Username = os.Getenv("DB_USER")
	ds.Password = os.Getenv("DB_PASSWORD")
	ds.Database = os.Getenv("DB_NAME")
}

func (ds *DataSource) CreateDB(name string) {

	db, err := sql.Open("postgres", ds.ShortDSN())
	if err != nil {
		log.Error(err)
	}
	defer db.Close()
	log.Info(ds.DSN())

	_, err = db.Exec("CREATE DATABASE " + name)
	if err != nil {
		log.Error(err)
	}
}

// TODO: Move common functionality between CreateDB and DropDB to separate private method
func (ds *DataSource) DopDB(name string) {
	db, err := sql.Open("postgres", ds.ShortDSN())
	if err != nil {
		log.Error(err)
	}
	defer db.Close()

	_, err = db.Exec("DROP DATABASE IF NOT EXISTS " + name)
	if err != nil {
		log.Error(err)
	}
}