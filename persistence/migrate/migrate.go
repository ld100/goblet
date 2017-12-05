package migrate

import (
	"time"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"

	"github.com/ld100/goblet/util/log"
	"github.com/ld100/goblet/persistence"
	"github.com/ld100/goblet/domain/users"
)

func Migrate() {
	db := persistence.GormDB

	db.LogMode(true)

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// create persons table
		{
			ID: "201712051642",
			Migrate: func(tx *gorm.DB) error {
				// it's a good pratice to copy the struct inside the function,
				// so side effects are prevented if the original struct changes during the time
				type User struct {
					ID        uint   `gorm:"primary_key"`
					CreatedAt time.Time
					UpdatedAt time.Time
					FirstName string `gorm:"size:255" valid:"optional"`              // Default size for string is 255, reset it with this tag
					LastName  string `gorm:"size:255" valid:"optional"`              // Default size for string is 255, reset it with this tag
					Email     string `gorm:"not null;unique" valid:"required,email"` // Set field as not nullable and unique
					Password  string `gorm:"size:255"`
				}
				return tx.AutoMigrate(&User{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("users").Error
			},
		},
	})

	if err := m.Migrate(); err != nil {
		log.Fatal("could not migrate: ", err)
	}
	log.Info("Migration did run successfully")
}

// Database Seed
// TODO: Wrap in transactions
func Seed() {
	log.Debug("Database seed initiated")
	user := users.User{
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
