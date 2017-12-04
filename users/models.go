package users

import (
	"fmt"
	"os"
	"time"
	"errors"
	"log"
	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/ld100/goblet/util/database"
)

// TODO: Remove this method, use full-fledged migrations and GORM connection polling instead
func MigrateUsers() {
	// Create database if not exist
	database.CreateDB(os.Getenv("DB_NAME"))

	connString := fmt.Sprintf(
		"host=%v user=%v dbname=%v sslmode=disable password=%v",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)
	db, err := gorm.Open("postgres", connString)

	if err != nil {
		panic(err)
		log.Fatal(err)
	}

	db.DropTable(&User{})
	db.AutoMigrate(&User{})

	password, err := HashPassword("password")
	if err != nil {
		panic(err)
		//log.Fatal(err)
	}

	user := User{
		FirstName:      "John",
		LastName:       "Doe",
		Email:          "you@example.com",
		PasswordDigest: password,
	}
	if db.NewRecord(user) {
		db.Create(&user)
	}

	defer db.Close()
}

type User struct {
	ID             uint   `gorm:"primary_key"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	FirstName      string `gorm:"size:255"`        // Default size for string is 255, reset it with this tag
	LastName       string `gorm:"size:255"`        // Default size for string is 255, reset it with this tag
	Email          string `gorm:"not null;unique"` // Set field as not nullable and unique
	PasswordDigest string
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
