package env

import (
	"github.com/ld100/goblet/pkg/persistence"
	"github.com/ld100/goblet/pkg/util/config"
	"github.com/ld100/goblet/pkg/util/logger"
)

// Environment context, which is passed through whole application.
// Contains Configuration, Database handler and Logger links
// Acts as a unit-tests-friendly alternative to global vars
type Env struct {
	DB     *persistence.DB
	Config *config.Config
	Logger *logger.Logger
}
