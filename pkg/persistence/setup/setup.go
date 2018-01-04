package setup

import (
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"

	"github.com/ld100/goblet/pkg/persistence"
	"github.com/ld100/goblet/pkg/persistence/migrate"
)

func SetupDatabases() {
	// Create database if not exist
	ds := &persistence.DataSource{}
	// Fetch database credentials from ENVIRONMENT
	ds.FetchENV()
	ds.CreateDB(os.Getenv("DB_NAME"))

	// Initiate global ORM var
	persistence.InitGormDB(ds)

	// Run migrations
	migrate.Migrate()

	// Run db seed
	migrate.Seed()
}
