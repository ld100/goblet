package config

import (
	"time"

	"github.com/spf13/viper"
)

// Initialize Viper configuration tool
func init() {
	viper.AutomaticEnv()
}

// Configuration wrapper on top of simple environment variables or spf13/viper
type Config struct {
}

func (cfg *Config) Get(key string) (interface{}) {
	return viper.Get(key)
}

func (cfg *Config) GetBool(key string) (bool) {
	return viper.GetBool(key)
}

func (cfg *Config) GetFloat64(key string) (float64) {
	return viper.GetFloat64(key)
}

func (cfg *Config) GetInt(key string) (int) {
	return viper.GetInt(key)
}

func (cfg *Config) GetString(key string) (string) {
	return viper.GetString(key)
}

func (cfg *Config) GetStringMap(key string) (map[string]interface{}) {
	return viper.GetStringMap(key)
}

func (cfg *Config) GetStringMapString(key string) (map[string]string) {
	return viper.GetStringMapString(key)
}

func (cfg *Config) GetStringSlice(key string) ([]string) {
	return viper.GetStringSlice(key)
}

func (cfg *Config) GetTime(key string) (time.Time) {
	return viper.GetTime(key)
}

func (cfg *Config) GetDuration(key string) (time.Duration) {
	return viper.GetDuration(key)
}

func (cfg *Config) IsSet(key string) (bool) {
	return viper.IsSet(key)
}

func (cfg *Config) Set(key string, value interface{}) {
	viper.Set(key, value)
}
