package migrate

import (
	"time"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"

	"github.com/ld100/goblet/pkg/persistence"
	"github.com/ld100/goblet/pkg/server/env"
)

func MigrateDB(env *env.Env) (err error) {
	cfg := env.Config
	log := env.Logger

	ds := persistence.NewDSFromCFG(cfg)
	handler, err := persistence.NewDB(ds)
	if err != nil {
		return err
	}
	db, err := handler.ORMConnection()
	if err != nil {
		return err
	}

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
		return err
	}
	log.Info("migrations did run successfully")

	return nil
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
