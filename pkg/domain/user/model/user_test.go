package model

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/ld100/goblet/pkg/util/securerandom"
	"github.com/stretchr/testify/assert"
)

func TestUserHashPassword(t *testing.T) {
	assert := assert.New(t)
	password := "12345zOMFG"
	hash, _ := HashPassword(password)
	assert.NotEqual(password, hash, "Password should have hash not equal to the password itself")
}

func TestUserCheckPasswordHash(t *testing.T) {
	assert := assert.New(t)
	password := "12345zOMFG"
	hash, _ := HashPassword(password)
	assert.Equal(CheckPasswordHash(password, hash), true, "Password %s hash %s is incorrect")
}

func TestUserBeforeUpdate(t *testing.T) {
	assert := assert.New(t)

	uuid, _ := securerandom.Uuid()
	password := "12345zOMFG"
	user := &User{
		ID:       1,
		Uuid:     uuid,
		Email:    "test@example.com",
		Password: password,
	}

	err := user.BeforeUpdate()
	if assert.Nil(err, "Password autohashing does not work") {
		if assert.NotEqual(user.Password, password, "Password was not hashed") {
			// Password was autohashed successfully
			// running BeforeUpdate for the second time, it should not hash password again
			tempHash := user.Password
			err := user.BeforeUpdate()
			if assert.Nil(err, "Password autohashing does not work") {
				assert.Equal(tempHash, user.Password, "Password was rehashed while not needed")
			}
		}
	}
}

func TestUserBeforeCreate(t *testing.T) {
	assert := assert.New(t)

	uuid, _ := securerandom.Uuid()
	password := "12345zOMFG"
	user := &User{
		ID:       1,
		Uuid:     uuid,
		Email:    "test@example.com",
		Password: password,
	}
	err := user.BeforeCreate()
	assert.Nil(err)
	assert.NotEqual(password, user.Password, "Password should be hashed on create")
}

func TestUserValidate(t *testing.T) {
	assert := assert.New(t)

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create sqlmock: %s", err)
	}

	gormDB, gerr := gorm.Open("postgres", db)
	if gerr != nil {
		t.Fatalf("can't open gorm connection: %s", err)
	}
	gormDB.LogMode(false)
	defer gormDB.Close()

	uuid, _ := securerandom.Uuid()
	password := "12345zOMFG"
	user := &User{
		ID:       1,
		Uuid:     uuid,
		Email:    "test@example.com",
		Password: password,
	}

	user.Validate(gormDB)
	assert.Equal(0, len(gormDB.GetErrors()))

	// Joe is forbidden name... just for the sake of example
	user.FirstName = "Joe"
	user.Validate(gormDB)
	assert.Equal(1, len(gormDB.GetErrors()))
}
