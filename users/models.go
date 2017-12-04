package users

import (
	"time"
	"errors"
	"log"
	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	"github.com/ld100/goblet/environment"
)

var db *gorm.DB

// TODO: Remove this test func
func MigrateUsers() {
	db := environment.GDB

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
	u.Password, err = HashPassword(u.Password)
	if err != nil {
		err = errors.New("cannot hash user password")
		log.Fatal(err)
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