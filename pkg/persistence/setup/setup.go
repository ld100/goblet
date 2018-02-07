package setup

import (
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"

	"github.com/ld100/goblet/pkg/persistence"
	"github.com/ld100/goblet/pkg/persistence/migrate"
)

func SetupDatabases() (*persistence.DB, error) {
	// Create database if not exist
	ds := &persistence.DataSource{}
	// Fetch database credentials from ENVIRONMENT
	ds.FetchENV()

	u, err := persistence.NewDButil(ds)
	if err != nil {
		return nil, err
	} else {
		err := u.CreateDB(os.Getenv("DB_NAME"))
		if err != nil {
			return nil, err
		}
	}

	db, err := persistence.NewDB(ds)

	// Initiate global ORM var
	// Used for backward-compatibility
	persistence.InitGormDB(ds)

	// Run migrations
	migrate.Migrate()

	// Run db seed
	migrate.Seed()

	return db, nil
}
