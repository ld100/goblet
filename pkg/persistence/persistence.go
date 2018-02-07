package persistence

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
	"github.com/qor/validations"

	"github.com/ld100/goblet/pkg/util/log"
)

// DB acts as a handler for databases, providing both Gorm and plain SQL connection interfaces
type DB struct {
	SqlDB  *sql.DB
	GormDB *gorm.DB
}

// GORM connection handler
func (db *DB) ORMConnection() (*gorm.DB, error) {
	if db.GormDB == nil {
		return nil, errors.New("ORM handler is not available")
	}
	return db.GormDB, nil
}

// SQL connection handler
func (db *DB) SQLConnection() (*sql.DB, error) {
	if db.SqlDB == nil {
		if db.GormDB == nil {
			return nil, errors.New("SQL handler is not available")
		}
		db.SqlDB = db.GormDB.DB()
	}
	return db.SqlDB, nil
}

// Constructor for database handlers, that provide both ORM/SQL connections
func NewDB(ds *DataSource) (*DB, error) {
	var err error
	db := &DB{}

	db.GormDB, err = gorm.Open("postgres", ds.DSN())
	if err != nil {
		log.Error(err)
		return nil, err
	}

	validations.RegisterCallbacks(db.GormDB)
	return db, nil
}

// Common database utils: create/drop database, etc
// Acts as a content-aware proxy adapter for sql.DB
type DBUtil struct {
	*sql.DB
}

// Create database with specified name
func (db *DBUtil) CreateDB(name string) (error) {
	_, err := db.DB.Exec("CREATE DATABASE " + name)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// Drop database with specified name if it exists
func (db *DBUtil) DropDB(name string) (error) {
	_, err := db.Exec("DROP DATABASE IF NOT EXISTS " + name)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func NewDButil(ds *DataSource) (*DBUtil, error) {
	var err error
	db := &DBUtil{}

	db.DB, err = sql.Open("postgres", ds.ShortDSN())
	//defer db.DB.Close()

	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err = db.DB.Ping(); err != nil {
		log.Error(err)
		return nil, err
	}

	return db, nil
}

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
	ds.Port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	ds.Username = os.Getenv("DB_USER")
	ds.Password = os.Getenv("DB_PASSWORD")
	ds.Database = os.Getenv("DB_NAME")
}

// =====================================================================
//  New implementation is on top of the file, old one is at the bottom
// =====================================================================

var GormDB *gorm.DB

func InitGormDB(ds *DataSource) {
	var err error
	GormDB, err = gorm.Open("postgres", ds.DSN())
	if err != nil {
		log.Error(err)
	}

	validations.RegisterCallbacks(GormDB)
}
