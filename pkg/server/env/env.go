package env

import (
	"github.com/ld100/goblet/pkg/persistence"
	//"github.com/ld100/goblet/pkg/util/log"
)

// Environment context, which is passed through whole application.
// Contains Configuration, Database handler and Logger links
// Acts as a unit-tests-friendly alternative to global vars
type Env struct {
	DB     *persistence.DB
	Config *Config
	Logger *Logger
}

// Configuration wrapper on top of simple environment variables or spf13/viper
type Config struct {
	Debug bool
}

// Logging wrapper on top of sirupsen/logrus
type Logger struct {
}
