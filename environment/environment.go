package environment

import (
	"log"
	"database/sql"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/qor/validations"
	_ "github.com/lib/pq"
)

var DB *sql.DB
var GDB *gorm.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = DB.Ping(); err != nil {
		log.Panic(err)
	}
}

func InitGDB(dataSourceName string) {
	var err error
	GDB, err = gorm.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	validations.RegisterCallbacks(GDB)
}
