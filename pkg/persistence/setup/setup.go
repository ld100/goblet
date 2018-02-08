package setup

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"

	"github.com/ld100/goblet/pkg/persistence"
	"github.com/ld100/goblet/pkg/util/config"
)

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

// TODO: Implement MigrateDB

// TODO: Implement SeedDB

// ======================================================
//  Old functionality, remove it once not used
// ======================================================

func SetupDatabases(cfg *config.Config) (*persistence.DB, error) {
	// Fetch database credentials from provided config
	ds := persistence.NewDSFromCFG(cfg)

	// Instantiate database handler
	db, err := persistence.NewDB(ds)
	if err != nil {
		return nil, err
	}

	return db, nil
}
