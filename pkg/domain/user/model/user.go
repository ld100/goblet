package model

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	"github.com/ld100/goblet/pkg/util/hash"
	"github.com/ld100/goblet/pkg/util/securerandom"
)

type User struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Uuid      string    `gorm:"not null;unique" json:"uuid"` // Set field as not nullable and unique
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	FirstName string    `gorm:"size:255" valid:"optional" json:"firstName"`          // Default size for string is 255, reset it with this tag
	LastName  string    `gorm:"size:255" valid:"optional" json:"lastName"`           // Default size for string is 255, reset it with this tag
	Email     string    `gorm:"not null;unique" valid:"required,email" json:"email"` // Set field as not nullable and unique
	Password  string    `gorm:"size:255" json:"password"`
	Sessions  []Session `json:"sessions"`
}

// GORM callback: Encode password before create
func (u *User) BeforeCreate() (err error) {
	// Hash password
	u.Password, err = HashPassword(u.Password)
	if err != nil {
		return errors.New("cannot hash user password: " + err.Error())
	}

	// Set UUID for the user
	u.Uuid, err = securerandom.Uuid()
	if err != nil {
		return errors.New("cannot generate UUID for user: " + err.Error())
	}

	return nil
}

// Detect if password was set, encode it if needed
func (u *User) BeforeUpdate() (err error) {
	if !hash.IsBcrypt(u.Password) {
		u.Password, err = HashPassword(u.Password)
		if err != nil {
			return errors.New("cannot hash user password: " + err.Error())
		}
	}
	return nil
}

func (u User) Validate(db *gorm.DB) {
	if u.FirstName == "Joe" {
		db.AddError(errors.New("go change your name, it is invalid"))
	}
}

func HashPassword(password string) (hash string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
