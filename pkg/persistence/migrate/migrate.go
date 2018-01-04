package migrate

import (
	"time"

	"github.com/jinzhu/gorm"
	gormigrate "gopkg.in/gormigrate.v1"

	"github.com/ld100/goblet/pkg/domain/user/model"
	"github.com/ld100/goblet/pkg/persistence"
	"github.com/ld100/goblet/pkg/util/log"
)

func Migrate() {
	db := persistence.GormDB

	db.LogMode(true)

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// create persons table
		{
			ID:       "201712051642",
			Migrate:  Migrate201712051642,
			Rollback: Rollback201712051642,
		},
		{
			ID:       "201712141900",
			Migrate:  Migrate201712141900,
			Rollback: Rollback201712141900,
		},
	})

	if err := m.Migrate(); err != nil {
		log.Fatal("could not migrate: ", err)
	}
	log.Info("Migration did run successfully")
}

func Migrate201712051642(tx *gorm.DB) error {
	// it's a good pratice to copy the struct inside the function,
	// so side effects are prevented if the original struct changes during the time
	type User struct {
		ID        uint   `gorm:"primary_key"`
		Uuid      string `gorm:"not null;unique"`
		CreatedAt time.Time
		UpdatedAt time.Time
		FirstName string `gorm:"size:255" valid:"optional"`              // Default size for string is 255, reset it with this tag
		LastName  string `gorm:"size:255" valid:"optional"`              // Default size for string is 255, reset it with this tag
		Email     string `gorm:"not null;unique" valid:"required,email"` // Set field as not nullable and unique
		Password  string `gorm:"size:255"`
	}
	return tx.AutoMigrate(&User{}).Error
}

func Rollback201712051642(tx *gorm.DB) error {
	return tx.DropTable("user").Error
}

func Migrate201712141900(tx *gorm.DB) error {
	type Session struct {
		ID        uint      `gorm:"primary_key"`
		Uuid      string    `gorm:"not null;unique"`
		UserID    uint      `gorm:"index"`
		ExpiresAt time.Time `json:"expiresAt"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}
	return tx.AutoMigrate(&Session{}).Error
}

func Rollback201712141900(tx *gorm.DB) error {
	return tx.DropTable("sessions").Error
}

// Database Seed
// TODO: Wrap in transactions
func Seed() {
	//log.Debug("Database seed initiated")
	user := model.User{
		FirstName: "Admin",
		LastName:  "Adminovich",
		Email:     "robot@example.com",
		Password:  "password",
	}

	db := persistence.GormDB
	error := db.Where("email = ?", user.Email).First(&user)
	if error != nil {
		// User with this e-mail was not found, so let's create one
		errors := db.Create(&user).GetErrors()
		if len(errors) > 0 {
			log.Fatal(errors)
		} else {
			log.Debug("created user entity with ID:", user.ID)
		}
	}
}
