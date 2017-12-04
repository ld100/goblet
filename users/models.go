package users

import (
	"fmt"
	"os"
	"time"
	"errors"
	"log"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/qor/validations"

	"github.com/ld100/goblet/util/database"
)

var db *gorm.DB

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
	validations.RegisterCallbacks(db)

	db.DropTable(&User{})
	db.AutoMigrate(&User{})

	user := User{
		FirstName: "John the",
		LastName:  "Doe",
		Email:     "you@example.com",
		Password:  "password",
	}
	if db.NewRecord(user) {
		//errs := db.Create(&user).GetErrors()
		//fmt.Print(errs)
		db.Create(&user)
	}

	defer db.Close()
}

type User struct {
	ID        uint   `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	FirstName string `gorm:"size:255" valid:"optional"`              // Default size for string is 255, reset it with this tag
	LastName  string `gorm:"size:255" valid:"optional"`              // Default size for string is 255, reset it with this tag
	Email     string `gorm:"not null;unique" valid:"required,email"` // Set field as not nullable and unique
	Password  string `gorm:"size:255"`
}

// GORM callback: Encode password before create
func (u *User) BeforeCreate() (err error) {
	// Do not encode already encoded password
	if !IsBase64(u.Password) {
		u.Password, err = HashPassword(u.Password)
		if err != nil {
			err = errors.New("cannot hash user password")
			log.Fatal(err)
		}
	}
	return
}

func (u User) Validate(db *gorm.DB) {
	if u.FirstName == "Joe" {
		db.AddError(errors.New("go change your name, it is invalid"))
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func IsBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}
