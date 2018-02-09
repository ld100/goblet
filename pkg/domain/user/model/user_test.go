package model

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/ld100/goblet/pkg/util/securerandom"
	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	password := "12345zOMFG"
	hash, _ := HashPassword(password)

	if password == hash {
		t.Errorf("Password %s should have hash not equal to the password itself", password)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "12345zOMFG"
	hash, _ := HashPassword(password)
	equals := CheckPasswordHash(password, hash)

	if !equals {
		t.Errorf("Password %s hash %s is incorrect", password, hash)
	}
}

func TestUserBeforeUpdate(t *testing.T) {
	uuid, _ := securerandom.Uuid()
	password := "12345zOMFG"
	user := &User{
		ID:       1,
		Uuid:     uuid,
		Email:    "test@example.com",
		Password: password,
	}

	err := user.BeforeUpdate()
	if err != nil {
		t.Errorf("Password autohashing does not work")
	} else {
		if user.Password == password {
			t.Errorf("Password autohashing does not work, password was not hashed")
		} else {
			// Password was autohashed successfully
			// running BeforeUpdate for the second time, it should not hash password again
			tempHash := user.Password
			err := user.BeforeUpdate()
			if err != nil {
				t.Errorf("Password autohashing does not work")
			} else {
				if tempHash != user.Password {
					t.Errorf("Password was rehashed while not needed")
				}
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
	err := user.BeforeUpdate()
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
