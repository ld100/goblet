package setup

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"

	"github.com/ld100/goblet/pkg/persistence"
	"github.com/ld100/goblet/pkg/persistence/migrate"
	"github.com/ld100/goblet/pkg/util/config"
)

func SetupDatabases(cfg *config.Config) (*persistence.DB, error) {
	// Create database if not exist
	ds := &persistence.DataSource{}
	// Fetch database credentials from ENVIRONMENT
	ds.FetchENV()

	u, err := persistence.NewDButil(ds)
	if err != nil {
		return nil, err
	} else {
		err := u.CreateDB(cfg.GetString("DB_NAME"))
		if err != nil {
			return nil, err
		}
	}

	db, err := persistence.NewDB(ds)
	if err != nil {
		return nil, err
	}

	// Run migrations
	migrate.Migrate(db)

	// Run db seed
	migrate.Seed(db)

	return db, nil
}
