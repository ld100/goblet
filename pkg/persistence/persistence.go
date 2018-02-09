package persistence

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
	"github.com/qor/validations"

	"github.com/ld100/goblet/pkg/util/config"
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

// Checks if database with specific name exists
func (db *DBUtil) Exists(name string) (bool, error) {
	query := fmt.Sprintf("SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = lower('%s');", name)
	rows, err := db.DB.Query(query)

	if err != nil {
		return false, err
	}

	defer rows.Close()
	var datname string;
	for rows.Next() {
		err = rows.Scan(&datname)
		if err != nil {
			return false, nil
		} else {
			return true, nil
		}
	}
	return false, nil
}

// Create database with specified name
func (db *DBUtil) CreateDB(name string) (error) {
	exists, err := db.Exists(name)
	if !exists {
		_, err = db.DB.Exec("CREATE DATABASE " + name)
		if err != nil {
			return err
		}
	}

	return nil
}

// Drop database with specified name if it exists
func (db *DBUtil) DropDB(name string) (error) {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
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
		return nil, err
	}

	if err = db.DB.Ping(); err != nil {
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

func NewDS(host string, port int, username string, password string, database string) (*DataSource) {
	ds := &DataSource{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
	if database != "" {
		ds.Database = database
	}
	return ds
}

func NewDSFromCFG(cfg *config.Config) (*DataSource) {
	ds := &DataSource{
		Host:     cfg.GetString("DB_HOST"),
		Port:     cfg.GetInt("DB_PORT"),
		Username: cfg.GetString("DB_USER"),
		Password: cfg.GetString("DB_PASSWORD"),
		Database: cfg.GetString("DB_NAME"),
	}

	return ds
}
