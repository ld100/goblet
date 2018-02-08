package setup

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"

	"github.com/ld100/goblet/pkg/persistence"
	"github.com/ld100/goblet/pkg/util/config"
)

// TODO: Fully move CreateDB functionality to command
// Creates SQL database with specified name
// Takes name from DB_NAME environment variable if not provided directly
func CreateDB(cfg *config.Config, name string) (err error) {
	ds := persistence.NewDSFromCFG(cfg)
	u, err := persistence.NewDButil(ds)
	if err != nil {
		return err
	}

	if len(name) == 0 {
		name = cfg.GetString("DB_NAME")
	}
	err = u.CreateDB(name)
	if err != nil {
		return err
	}
	return nil
}

// TODO: Implement DropDB
