package seed

import (
	"github.com/ld100/goblet/pkg/domain/user/model"
	"github.com/ld100/goblet/pkg/persistence"
	"github.com/ld100/goblet/pkg/server/env"
)

// Database Seed
// TODO: Wrap in transactions: http://gorm.io/advanced.html#transactions
func SeedDB(env *env.Env) (err error) {
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

	user := model.User{
		FirstName: "Admin",
		LastName:  "Adminovich",
		Email:     "robot@example.com",
		Password:  "password",
	}

	if err = db.Where("email = ?", user.Email).First(&user).Error; err != nil {
		// User with this e-mail was not found, so let's create one
		errors := db.Create(&user).GetErrors()
		if len(errors) > 0 {
			log.Fatal(errors)
		} else {
			log.Debug("created user entity with ID:", user.ID)
		}
	}
	return nil
}
