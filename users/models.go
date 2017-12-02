package models

import (
	"fmt"
	"os"
	"time"
	"errors"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	connString := fmt.Sprintf(
		"host=%s user=%s dbname=%s sslmode=disable password=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)
	db, err := gorm.Open(connString)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
}

type User struct {
	ID        uint   `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	FirstName string `gorm:"size:255"`        // Default size for string is 255, reset it with this tag
	LastName  string `gorm:"size:255"`        // Default size for string is 255, reset it with this tag
	Email     string `gorm:"not null;unique"` // Set field as not nullable and unique
}

// GORM callback
func (u *User) BeforeSave() (err error) {
	err = u.Validate()
	if err != nil {
		err = errors.New("read only user")
	}
	return
}

func (u *User) Validate() (err error) {
	fmt.Println("Validating User model before save")
	return
}
