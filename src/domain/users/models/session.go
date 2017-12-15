package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/ld100/goblet/util/log"
	"github.com/ld100/goblet/util/securerandom"
)

type Session struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Uuid      string    `gorm:"not null;unique" json:"uuid"` // Set field as not nullable and unique
	UserID    uint      `gorm:"index" valid:"required"`
	ExpiresAt time.Time `json:"expiresAt" valid:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (s Session) Validate(db *gorm.DB) {
	if s.ExpiresAt.Before(time.Now().UTC()) {
		db.AddError(errors.New("you cannot set session expiration tim in the past"))
	}
}

func (s *Session) BeforeCreate() (err error) {
	// Set UUID for the user
	if len(s.Uuid) == 0 {
		s.Uuid, err = securerandom.Uuid()
		if err != nil {
			err = errors.New("cannot generate UUID for session")
			log.Fatal(err)
		}
	}

	return
}

func CleanUpSessions() (err error) {
	// TODO: Implement method that cleans up all older sessions
	return nil
}