package model

import (
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestSessionBeforeCreate(t *testing.T) {
	password := "12345zOMFG"
	user := &User{
		ID:       1,
		Uuid:     "",
		Email:    "test@example.com",
		Password: password,
	}

	session := &Session{
		ID:        1,
		Uuid:      "",
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC(),
	}

	err := session.BeforeCreate()
	if err != nil {
		t.Errorf("UUID generation for sessions does not work: error")
	} else {
		if len(session.Uuid) < 1 {
			t.Errorf("UUID generation for sessions does not work: UUID is empty")
		}
	}
}

func TestSessionValidate(t *testing.T) {
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

	password := "12345zOMFG"
	user := &User{
		ID:       1,
		Uuid:     "",
		Email:    "test@example.com",
		Password: password,
	}

	session := &Session{
		ID:        1,
		Uuid:      "",
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC(),
	}

	session.Validate(gormDB)
	user.Validate(gormDB)
	assert.Equal(1, len(gormDB.GetErrors()), "Incorrect session ExpiresAt")
}

func TestSessionCleanUpSessions(t *testing.T) {
	assert := assert.New(t)
	err := CleanUpSessions()
	assert.Nil(err, "TestCleanUpSessions should not work, since it is not implemented yet")
}
