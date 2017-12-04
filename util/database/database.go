package database

import (
	"fmt"
	"os"
	"database/sql"

	_ "github.com/lib/pq"
)

var dataSourceName string = fmt.Sprintf(
	"host=%s user=%s sslmode=disable password=%s",
	os.Getenv("DB_HOST"),
	os.Getenv("DB_USER"),
	os.Getenv("DB_PASSWORD"),
)

func CreateDB(name string) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println(db)

	_, err = db.Exec("CREATE DATABASE " + name)
	if err != nil {
		panic(err)
	}
}

// TODO: Move common functionality between CreateDB and DropDB to separate private method
func DropDB(name string) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("DROP DATABASE IF NOT EXISTS " + name)
	if err != nil {
		panic(err)
	}
}
